package utils

import (
    "bytes"
    "database/sql"
    "encoding/binary"
    "errors"
    "fmt"
    "io"
    "math"
    "net"
    "reflect"
    "strings"
    "time"
)

const (
    SphVerMajorProto        = 0x1
    SphVerCommandSearch     = 0x119 // 0x11D for 2.1
    SphVerCommandExcerpt    = 0x104
    SphVerCommandUpdate     = 0x102 // 0x103 for 2.1
    SphVerCommandKeywords   = 0x100
    SphVerCommandStatus     = 0x100
    SphVerCommandFlushattrs = 0x100
)

const (
    SphMatchAll = iota
    SphMatchAny
    SphMatchPhrase
    SphMatchBoolean
    SphMatchExtended
    SphMatchFullscan
    SphMatchExtended2
)

const (
    SphRankProximityBm25 = iota
    SphRankBm25
    SphRankNone
    SphRankWordcount
    SphRankProximity
    SphRankMatchany
    SphRankFieldmask
    SphRankSph04
    SphRankExpr
    SphRankTotal
)

const (
    SphSortRelevance = iota
    SphSortAttrDesc
    SphSortAttrAsc
    SphSortTimeSegments
    SphSortExtended
    SphSortExpr
)

const (
    SphGroupbyDay = iota
    SphGroupbyWeek
    SphGroupbyMonth
    SphGroupbyYear
    SphGroupbyAttr
    SphGroupbyAttrpair
)

const (
    SearchdOk = iota
    SearchdError
    SearchdRetry
    SearchdWarning
)

const (
    SphAttrNone = iota
    SphAttrInteger
    SphAttrTimestamp
    SphAttrOrdinal
    SphAttrBool
    SphAttrFloat
    SphAttrBigint
    SphAttrString
    SphAttrMulti   = 0x40000001 // 0x40000000
    SphAttrMulti64 = 0x40000002
)

const (
    SearchdCommandSearch = iota
    SearchdCommandExcerpt
    SearchdCommandUpdate
    SearchdCommandKeywords
    SearchdCommandPersist
    SearchdCommandStatus
    SearchdCommandQuery
    SearchdCommandFlushattrs
)

const (
    SphFilterValues = iota
    SphFilterRange
    SphFilterFloatrange
)

type sphFilter struct {
    attr       string
    filterType int
    values     []uint64
    umin       uint64
    umax       uint64
    fmin       float32
    fmax       float32
    exclude    bool
}

type sphOverride struct {
    attrName string
    attrType int
    values   map[uint64]interface{}
}

type sphinxMatch struct {
    DocId      uint64                 `json:"doc_id"`
    Weight     int                    `json:"weight"`
    AttrValues map[string]interface{} `json:"attr_values"`
}

type sphinxWordInfo struct {
    Word string `json:"word"`
    Docs int    `json:"docs"`
    Hits int    `json:"hits"`
}

type SphinxResult struct {
    Fields     []string         `json:"fields"`     // full-text field namess.
    AttrNames  []string         `json:"-"`          // attribute names.
    AttrTypes  []int            `json:"-"`          // attribute types.
    Attrs      map[string]int   `json:"attrs"`      // attribute names -> types.
    Matches    []sphinxMatch    `json:"matches"`    // retrieved matches.
    Words      []sphinxWordInfo `json:"words"`      // per-word statistics.
    Total      int              `json:"total"`      // total matches in this result set.
    TotalFound int              `json:"totalFound"` // total matches found in the index(es).
    Time       float32          `json:"time"`       // elapsed time (as reported by searchd), in seconds.
    Warning    string           `json:"warning"`    // last warning message.
    Error      error            `json:"error"`      // last error.
    Status     int              `json:"status"`     // query status (refer to SEARCHD_xxx constants in Client).
}

type SphOptions struct {
    Host                     string
    Port                     int
    Socket                   string
    SqlPort                  int
    SqlSocket                string
    RetryCount               int
    RetryDelay               int
    Timeout                  int
    Offset                   int
    Limit                    int
    MaxMatches               int
    Cutoff                   int
    MaxQueryTimeMilliseconds int
    Select                   string
    MatchMode                int
    RankMode                 int
    RankExpr                 string
    SortMode                 int
    SortBy                   string // attribute to sort by (defualt is "")
    MinId                    uint64 // min ID to match (default is 0, which means no limit)
    MaxId                    uint64 // max ID to match (default is 0, which means no limit)
    LatitudeAttr             string
    LongitudeAttr            string
    Latitude                 float32
    Longitude                float32
    GroupBy                  string // group-by attribute name
    GroupFunc                int    // group-by function (to pre-process group-by attribute value with)
    GroupSort                string // group-by sorting clause (to sort groups in result set with)
    GroupDistinct            string // group-by count-distinct attribute

    // for sphinxql
    Index   string
    Columns []string
    Where   string
}

type SphinxClient struct {
    *SphOptions
    conn         net.Conn
    warning      string
    err          error
    connerror    bool
    weights      []int
    filters      []sphFilter
    reqs         [][]byte
    indexWeights map[string]int
    fieldWeights map[string]int
    overrides    map[string]sphOverride

    // For sphinxql
    DB  *sql.DB
    val reflect.Value
}

// defaultSphinxOptions 默认配置
var defaultSphinxOptions = &SphOptions{
    Host:       "localhost",
    Port:       3312,
    SqlPort:    9306,
    Offset:     0,
    Limit:      20,
    MaxMatches: 1000,
    MatchMode:  SphMatchAll,
    SortMode:   SphSortRelevance,
    RankMode:   SphRankProximityBm25,
    GroupFunc:  SphGroupbyDay,
    GroupSort:  "@group desc",
    Timeout:    2000,
    Select:     "*",
}

func NewSphinxClient(opts ...*SphOptions) (sc *SphinxClient) {
    if opts != nil {
        return &SphinxClient{SphOptions: opts[0]}
    }
    return &SphinxClient{
        SphOptions:   defaultSphinxOptions,
        weights:      make([]int, 0),
        filters:      make([]sphFilter, 0),
        indexWeights: make(map[string]int),
        fieldWeights: make(map[string]int),
        overrides:    make(map[string]sphOverride),
    }
}

func (sc *SphinxClient) GetLastError() error {
    return sc.err
}

func (sc *SphinxClient) Error() error {
    return sc.err
}

func (sc *SphinxClient) GetLastWarning() string {
    return sc.warning
}

func (sc *SphinxClient) SetServer(host string, port int) *SphinxClient {
    isTcpMode := true

    if host != "" {
        if host[0] == '/' {
            sc.Socket = host
            isTcpMode = false
        } else if len(host) > 7 && host[:7] == "unix://" {
            sc.Socket = host[7:]
            isTcpMode = false
        } else {
            sc.Host = host
        }
    } else {
        sc.Host = defaultSphinxOptions.Host
    }

    if isTcpMode {
        if port > 0 {
            sc.Port = port
        } else {
            sc.Port = defaultSphinxOptions.Port
        }
    }

    return sc
}

func (sc *SphinxClient) SetSqlServer(host string, sqlport int) *SphinxClient {
    isTcpMode := true

    if host != "" {
        if host[0] == '/' {
            sc.SqlSocket = host
            isTcpMode = false
        } else if len(host) > 7 && host[:7] == "unix://" {
            sc.SqlSocket = host[7:]
            isTcpMode = false
        } else {
            sc.Host = host
        }
    } else {
        sc.Host = defaultSphinxOptions.Host
    }

    if isTcpMode {
        if sqlport > 0 {
            sc.SqlPort = sqlport
        } else {
            sc.SqlPort = defaultSphinxOptions.SqlPort
        }
    }

    return sc
}

func (sc *SphinxClient) SetRetries(count, delay int) *SphinxClient {
    if count < 0 {
        sc.err = fmt.Errorf("SetRetries > count must not be negative: %d", count)
        return sc
    }
    if delay < 0 {
        sc.err = fmt.Errorf("SetRetries > delay must not be negative: %d", delay)
        return sc
    }

    sc.RetryCount = count
    sc.RetryDelay = delay
    return sc
}

func (sc *SphinxClient) SetConnectTimeout(millisecond int) *SphinxClient {
    if millisecond < 0 {
        sc.err = fmt.Errorf("SetConnectTimeout > connect timeout must not be negative: %d", millisecond)
        return sc
    }

    sc.Timeout = millisecond
    return sc
}

func (sc *SphinxClient) IsConnectError() bool {
    return sc.connerror
}

func (sc *SphinxClient) SetLimits(offset, limit, maxMatches, cutoff int) *SphinxClient {
    if offset < 0 {
        sc.err = fmt.Errorf("SetLimits > offset must not be negative: %d", offset)
        return sc
    }
    if limit <= 0 {
        sc.err = fmt.Errorf("SetLimits > limit must be positive: %d", limit)
        return sc
    }
    if maxMatches <= 0 {
        sc.err = fmt.Errorf("SetLimits > maxMatches must be positive: %d", maxMatches)
        return sc
    }
    if cutoff < 0 {
        sc.err = fmt.Errorf("SetLimits > cutoff must not be negative: %d", cutoff)
        return sc
    }

    sc.Offset = offset
    sc.Limit = limit
    sc.MaxMatches = maxMatches
    if cutoff > 0 {
        sc.Cutoff = cutoff
    }
    return sc
}

func (sc *SphinxClient) SetMaxQueryTime(milliseconds int) *SphinxClient {
    if milliseconds < 0 {
        sc.err = fmt.Errorf("SetMaxQueryTime > maxQueryTime must not be negative: %d", milliseconds)
        return sc
    }

    sc.MaxQueryTimeMilliseconds = milliseconds
    return sc
}

func (sc *SphinxClient) SetOverride(attrName string, attrType int, values map[uint64]interface{}) *SphinxClient {
    if attrName == "" {
        sc.err = errors.New("SetOverride attrName is empty")
        return sc
    }
    if (attrType < 1 || attrType > SphAttrString) && attrType != SphAttrMulti && SphAttrMulti != SphAttrMulti64 {
        sc.err = fmt.Errorf("SetOverride > invalid attrType: %d", attrType)
        return sc
    }

    sc.overrides[attrName] = sphOverride{
        attrName: attrName,
        attrType: attrType,
        values:   values,
    }
    return sc
}

func (sc *SphinxClient) SetSelect(s string) *SphinxClient {
    if s == "" {
        sc.err = errors.New("SetSelect > selectStr is empty")
        return sc
    }

    sc.Select = s
    return sc
}

func (sc *SphinxClient) SetMatchMode(mode int) *SphinxClient {
    if mode < 0 || mode > SphMatchExtended2 {
        sc.err = fmt.Errorf("SetMatchMode > unknown mode value; use one of the SPH_MATCH_xxx constants: %d", mode)
        return sc
    }

    sc.MatchMode = mode
    return sc
}

func (sc *SphinxClient) SetRankingMode(ranker int, rankexpr ...string) *SphinxClient {
    if ranker < 0 || ranker > SphRankTotal {
        sc.err = fmt.Errorf("SetRankingMode > unknown ranker value; use one of the SPH_RANK_xxx constants: %d", ranker)
        return sc
    }

    sc.RankMode = ranker

    if len(rankexpr) > 0 {
        if ranker != SphRankExpr {
            sc.err = fmt.Errorf("SetRankingMode > rankexpr must used with SPH_RANK_EXPR! ranker: %d  rankexpr: %s", ranker, rankexpr)
            return sc
        }

        sc.RankExpr = rankexpr[0]
    }

    return sc
}

func (sc *SphinxClient) SetSortMode(mode int, sortBy string) *SphinxClient {
    if mode < 0 || mode > SphSortExpr {
        sc.err = fmt.Errorf("SetSortMode > unknown mode value; use one of the available SPH_SORT_xxx constants: %d", mode)
        return sc
    }
    // SPH_SORT_RELEVANCE ignores any additional parameters and always sorts matches by relevance rank.
    // All other modes require an additional sorting clause.
    if (mode != SphSortRelevance) && (sortBy == "") {
        sc.err = fmt.Errorf("SetSortMode > sortby string must not be empty in selected mode: %d", mode)
        return sc
    }

    sc.SortMode = mode
    sc.SortBy = sortBy
    return sc
}

func (sc *SphinxClient) SetFieldWeight(name string, val int) *SphinxClient {
    if val < 1 {
        sc.err = fmt.Errorf("SetFieldWeight > weight must be positive 32-bit integers, field:%s  weight:%d", name, val)
        return sc
    }
    if sc.fieldWeights == nil {
        sc.fieldWeights = map[string]int{
            name: val,
        }
    } else {
        sc.fieldWeights[name] = val
    }
    return sc
}

func (sc *SphinxClient) SetFieldWeights(weights map[string]int) *SphinxClient {
    // 默认权重值为1
    for field, weight := range weights {
        if weight < 1 {
            sc.err = fmt.Errorf("SetFieldWeights > weights must be positive 32-bit integers, field:%s  weight:%d", field, weight)
            return sc
        }
    }

    sc.fieldWeights = weights
    return sc
}

func (sc *SphinxClient) SetIndexWeight(name string, val int) *SphinxClient {
    if val < 1 {
        sc.err = fmt.Errorf("SetIndexWeight > weight must be positive 32-bit integers, field:%s  weight:%d", name, val)
        return sc
    }
    if sc.indexWeights == nil {
        sc.indexWeights = map[string]int{
            name: val,
        }
    } else {
        sc.indexWeights[name] = val
    }
    return sc
}

func (sc *SphinxClient) SetIndexWeights(weights map[string]int) *SphinxClient {
    // 默认权重值为1
    for field, weight := range weights {
        if weight < 1 {
            sc.err = fmt.Errorf("SetIndexWeights > weights must be positive 32-bit integers, field:%s  weight:%d", field, weight)
            return sc
        }
    }

    sc.indexWeights = weights
    return sc
}

func (sc *SphinxClient) SetIDRange(min, max uint64) *SphinxClient {
    if min > max {
        sc.err = fmt.Errorf("SetIDRange > min > max! min:%d  max:%d", min, max)
        return sc
    }

    sc.MinId = min
    sc.MaxId = max
    return sc
}

func (sc *SphinxClient) SetFilter(attr string, values []uint64, exclude bool) *SphinxClient {
    if attr == "" {
        sc.err = fmt.Errorf("SetFilter > attribute name is empty")
        return sc
    }
    if len(values) == 0 {
        sc.err = fmt.Errorf("SetFilter > values is empty")
        return sc
    }

    sc.filters = append(sc.filters, sphFilter{
        filterType: SphFilterValues,
        attr:       attr,
        values:     values,
        exclude:    exclude,
    })
    return sc
}

func (sc *SphinxClient) SetFilterRange(attr string, umin, umax uint64, exclude bool) *SphinxClient {
    if attr == "" {
        sc.err = fmt.Errorf("SetFilterRange > attribute name is empty")
        return sc
    }
    if umin > umax {
        sc.err = fmt.Errorf("SetFilterRange > min > max! umin:%d  umax:%d", umin, umax)
        return sc
    }

    sc.filters = append(sc.filters, sphFilter{
        filterType: SphFilterRange,
        attr:       attr,
        umin:       umin,
        umax:       umax,
        exclude:    exclude,
    })
    return sc
}

func (sc *SphinxClient) SetFilterFloatRange(attr string, fmin, fmax float32, exclude bool) *SphinxClient {
    if attr == "" {
        sc.err = fmt.Errorf("SetFilterFloatRange > attribute name is empty")
        return sc
    }
    if fmin > fmax {
        sc.err = fmt.Errorf("SetFilterFloatRange > min > max")
        return sc
    }

    sc.filters = append(sc.filters, sphFilter{
        filterType: SphFilterFloatrange,
        attr:       attr,
        fmin:       fmin,
        fmax:       fmax,
        exclude:    exclude,
    })
    return sc
}

func (sc *SphinxClient) SetGeoAnchor(latitudeAttr, longitudeAttr string, latitude, longitude float32) *SphinxClient {
    if latitudeAttr == "" {
        sc.err = fmt.Errorf("SetGeoAnchor > latitudeAttr is empty")
        return sc
    }
    if longitudeAttr == "" {
        sc.err = fmt.Errorf("SetGeoAnchor > longitudeAttr is empty")
        return sc
    }

    sc.LatitudeAttr = latitudeAttr
    sc.LongitudeAttr = longitudeAttr
    sc.Latitude = latitude
    sc.Longitude = longitude
    return sc
}

func (sc *SphinxClient) SetGroupBy(groupBy string, groupFunc int, groupSort string) *SphinxClient {
    if groupFunc < 0 || groupFunc > SphGroupbyAttrpair {
        sc.err = fmt.Errorf("SetGroupBy > unknown groupFunc value: '%d', use one of the available SPH_GROUPBY_xxx constants",
            groupFunc)
        return sc
    }

    sc.GroupBy = groupBy
    sc.GroupFunc = groupFunc
    sc.GroupSort = groupSort
    return sc
}

func (sc *SphinxClient) SetGroupDistinct(groupDistinct string) *SphinxClient {
    if groupDistinct == "" {
        sc.err = errors.New("SetGroupDistinct > groupDistinct is empty")
        return sc
    }
    sc.GroupDistinct = groupDistinct
    return sc
}

func (sc *SphinxClient) Query(query, index, comment string) (result *SphinxResult, err error) {
    if index == "" {
        index = "*"
    }

    // 重置请求数据
    sc.reqs = nil

    if _, err = sc.AddQuery(query, index, comment); err != nil {
        return nil, err
    }

    results, er := sc.RunQueries()
    if er != nil {
        return nil, er
    }
    if len(results) == 0 {
        return nil, fmt.Errorf("Query > Empty results!\nClient: %#v", sc)
    }

    result = &results[0]
    if result.Error != nil {
        return nil, fmt.Errorf("result error: %s", result.Error.Error())
    }

    sc.warning = result.Warning
    return
}

func (sc *SphinxClient) AddQuery(query, index, comment string) (i int, err error) {
    var req []byte

    req = writeInt32ToBytes(req, sc.Offset)
    req = writeInt32ToBytes(req, sc.Limit)
    req = writeInt32ToBytes(req, sc.MatchMode)
    req = writeInt32ToBytes(req, sc.RankMode)
    if sc.RankMode == SphRankExpr {
        req = writeLenStrToBytes(req, sc.RankExpr)
    }
    req = writeInt32ToBytes(req, sc.SortMode)
    req = writeLenStrToBytes(req, sc.SortBy)
    req = writeLenStrToBytes(req, query)
    req = writeInt32ToBytes(req, len(sc.weights))
    if sc.weights != nil {
        for _, w := range sc.weights {
            req = writeInt32ToBytes(req, w)
        }
    }
    req = writeLenStrToBytes(req, index)
    req = writeInt32ToBytes(req, 1) // id64 range marker
    req = writeInt64ToBytes(req, sc.MinId)
    req = writeInt64ToBytes(req, sc.MaxId)

    // filters
    req = writeInt32ToBytes(req, len(sc.filters))
    if sc.filters != nil {
        for _, f := range sc.filters {
            req = writeLenStrToBytes(req, f.attr)
            req = writeInt32ToBytes(req, f.filterType)

            switch f.filterType {
            case SphFilterValues:
                req = writeInt32ToBytes(req, len(f.values))
                for _, v := range f.values {
                    req = writeInt64ToBytes(req, v)
                }
            case SphFilterRange:
                req = writeInt64ToBytes(req, f.umin)
                req = writeInt64ToBytes(req, f.umax)
            case SphFilterFloatrange:
                req = writeFloat32ToBytes(req, f.fmin)
                req = writeFloat32ToBytes(req, f.fmax)
            }

            if f.exclude {
                req = writeInt32ToBytes(req, 1)
            } else {
                req = writeInt32ToBytes(req, 0)
            }
        }
    }

    req = writeInt32ToBytes(req, sc.GroupFunc)
    req = writeLenStrToBytes(req, sc.GroupBy)
    req = writeInt32ToBytes(req, sc.MaxMatches)
    req = writeLenStrToBytes(req, sc.GroupSort)
    req = writeInt32ToBytes(req, sc.Cutoff)
    req = writeInt32ToBytes(req, sc.RetryCount)
    req = writeInt32ToBytes(req, sc.RetryDelay)
    req = writeLenStrToBytes(req, sc.GroupDistinct)

    if sc.LatitudeAttr == "" || sc.LongitudeAttr == "" {
        req = writeInt32ToBytes(req, 0)
    } else {
        req = writeInt32ToBytes(req, 1)
        req = writeLenStrToBytes(req, sc.LatitudeAttr)
        req = writeLenStrToBytes(req, sc.LongitudeAttr)
        req = writeFloat32ToBytes(req, sc.Latitude)
        req = writeFloat32ToBytes(req, sc.Longitude)
    }

    // indexWeights
    req = writeInt32ToBytes(req, len(sc.indexWeights))
    if sc.indexWeights != nil {
        for ind, wei := range sc.indexWeights {
            req = writeLenStrToBytes(req, ind)
            req = writeInt32ToBytes(req, wei)
        }
    }

    // maxQueryTime
    req = writeInt32ToBytes(req, sc.MaxQueryTimeMilliseconds)

    // fieldWeights
    req = writeInt32ToBytes(req, len(sc.fieldWeights))
    if sc.fieldWeights != nil {
        for fie, wei := range sc.fieldWeights {
            req = writeLenStrToBytes(req, fie)
            req = writeInt32ToBytes(req, wei)
        }
    }

    // comment
    req = writeLenStrToBytes(req, comment)

    // attribute overrides
    req = writeInt32ToBytes(req, len(sc.overrides))
    if sc.overrides != nil {
        for _, override := range sc.overrides {
            req = writeLenStrToBytes(req, override.attrName)
            req = writeInt32ToBytes(req, override.attrType)
            req = writeInt32ToBytes(req, len(override.values))
            for id, v := range override.values {
                req = writeInt64ToBytes(req, id)
                switch override.attrType {
                case SphAttrInteger:
                    req = writeInt32ToBytes(req, v.(int))
                case SphAttrFloat:
                    req = writeFloat32ToBytes(req, v.(float32))
                case SphAttrBigint:
                    req = writeInt64ToBytes(req, v.(uint64))
                default:
                    return -1, fmt.Errorf("AddQuery > attr value is not int/float32/uint64")
                }
            }
        }
    }
    req = writeLenStrToBytes(req, sc.Select)

    // 单个查询压入请求体
    sc.reqs = append(sc.reqs, req)
    return len(sc.reqs) - 1, nil
}

func (sc *SphinxClient) RunQueries() (results []SphinxResult, err error) {
    if len(sc.reqs) == 0 {
        return nil, fmt.Errorf("RunQueries > No queries defined, issue AddQuery() first")
    }

    // 构造全部查询数据
    nreqs := len(sc.reqs)
    var allReqs []byte
    allReqs = writeInt32ToBytes(allReqs, 0)
    allReqs = writeInt32ToBytes(allReqs, nreqs)
    for _, req := range sc.reqs {
        allReqs = append(allReqs, req...)
    }
    // 执行请求
    response, er := sc.doRequest(SearchdCommandSearch, SphVerCommandSearch, allReqs)
    if er != nil {
        return nil, er
    }

    var bp = byteParser{stream: response}
    for i := 0; i < nreqs; i++ {
        var result = SphinxResult{Status: -1}
        result.Status = bp.Int32()
        if result.Status != SearchdOk {
            message := bp.String()
            if result.Status == SearchdWarning {
                result.Warning = message
            } else {
                result.Error = errors.New(message)
                results = append(results, result)
                continue
            }
        }

        // read schema
        nfields := bp.Int32()
        result.Fields = make([]string, nfields)
        for fieldNum := 0; fieldNum < nfields; fieldNum++ {
            result.Fields[fieldNum] = bp.String()
        }

        numAttrs := bp.Int32()
        result.AttrNames = make([]string, numAttrs)
        result.AttrTypes = make([]int, numAttrs)
        result.Attrs = make(map[string]int)
        for attrNum := 0; attrNum < numAttrs; attrNum++ {
            result.AttrNames[attrNum] = bp.String()
            result.AttrTypes[attrNum] = bp.Int32()
            result.Attrs[result.AttrNames[attrNum]] = result.AttrTypes[attrNum]
        }

        // read match count
        count := bp.Int32()
        id64 := bp.Int32()
        result.Matches = make([]sphinxMatch, count)
        for matchesNum := 0; matchesNum < count; matchesNum++ {
            var match sphinxMatch
            if id64 == 1 {
                match.DocId = bp.Uint64()
            } else {
                match.DocId = uint64(bp.Uint32())
            }
            match.Weight = bp.Int32()
            match.AttrValues = make(map[string]interface{})

            for attrNum := 0; attrNum < len(result.AttrTypes); attrNum++ {
                attrName := result.AttrNames[attrNum]
                attrType := result.AttrTypes[attrNum]
                switch attrType {
                case SphAttrBigint:
                    match.AttrValues[attrName] = bp.Uint64()
                case SphAttrFloat:
                    f, err := bp.Float32()
                    if err != nil {
                        return nil, fmt.Errorf("binary.Read error: %v", err)
                    }
                    match.AttrValues[attrName] = f
                case SphAttrString:
                    match.AttrValues[attrName] = bp.String()
                case SphAttrMulti:
                    nvals := bp.Int32()
                    var vals = make([]uint32, nvals)
                    for valNum := 0; valNum < nvals; valNum++ {
                        vals[valNum] = bp.Uint32()
                    }
                    match.AttrValues[attrName] = vals
                case SphAttrMulti64:
                    nvals := bp.Int32()
                    nvals = nvals / 2
                    var vals = make([]uint64, nvals)
                    for valNum := 0; valNum < nvals; valNum++ {
                        vals[valNum] = uint64(bp.Uint32())
                        bp.Uint32()
                    }
                    match.AttrValues[attrName] = vals
                default:
                    match.AttrValues[attrName] = bp.Uint32()
                }
            }
            result.Matches[matchesNum] = match
        }

        result.Total = bp.Int32()
        result.TotalFound = bp.Int32()
        result.Time = float32(bp.Uint32()) / 1000.0

        nwords := bp.Int32()
        result.Words = make([]sphinxWordInfo, nwords)
        for wordNum := 0; wordNum < nwords; wordNum++ {
            result.Words[wordNum].Word = bp.String()
            result.Words[wordNum].Docs = bp.Int32()
            result.Words[wordNum].Hits = bp.Int32()
        }

        results = append(results, result)
    }

    // 重置请求数据
    sc.reqs = nil
    return
}

func (sc *SphinxClient) ResetFilters() {
    sc.filters = []sphFilter{}
    sc.LatitudeAttr = ""
    sc.LongitudeAttr = ""
    sc.Latitude = 0.0
    sc.Longitude = 0.0
}

func (sc *SphinxClient) ResetGroupBy() {
    sc.GroupBy = ""
    sc.GroupFunc = SphGroupbyDay
    sc.GroupSort = "@group desc"
    sc.GroupDistinct = ""
}

// ExcerptsOpts 额外附带功能
type ExcerptsOpts struct {
    BeforeMatch        string // default is "<b>".
    AfterMatch         string // default is "</b>".
    ChunkSeparator     string // A string to insert between snippet chunks (passages). Default is " ... ".
    Limit              int    // Maximum snippet size, in symbols (codepoints). default is 256.
    Around             int    // How much words to pick around each matching keywords block. default is 5.
    ExactPhrase        bool   // Whether to highlight exact query phrase matches only instead of individual keywords.
    SinglePassage      bool   // Whether to extract single best passage only.
    UseBoundaries      bool   // Whether to additionaly break passages by phrase boundary characters, as configured in index settings with phrase_boundary directive.
    WeightOrder        bool   // Whether to sort the extracted passages in order of relevance (decreasing weight), or in order of appearance in the document (increasing position).
    QueryMode          bool   // Whether to handle $words as a query in extended syntax, or as a bag of words (default behavior).
    ForceAllWords      bool   // Ignores the snippet length limit until it includes all the keywords.
    LimitPassages      int    // Limits the maximum number of passages that can be included into the snippet. default is 0 (no limit).
    LimitWords         int    // Limits the maximum number of keywords that can be included into the snippet. default is 0 (no limit).
    StartPassageId     int    // Specifies the starting value of %PASSAGE_ID% macro (that gets detected and expanded in BeforeMatch, AfterMatch strings). default is 1.
    LoadFiles          bool   // Whether to handle $docs as data to extract snippets from (default behavior), or to treat it as file names, and load data from specified files on the server side.
    LoadFilesScattered bool   // It assumes "load_files" option, and works only with distributed snippets generation with remote agents. The source files for snippets could be distributed among different agents, and the main daemon will merge together all non-erroneous results. So, if one agent of the distributed index has 'file1.txt', another has 'file2.txt' and you call for the snippets with both these files, the sphinx will merge results from the agents together, so you will get the snippets from both 'file1.txt' and 'file2.txt'.
    HtmlStripMode      string // HTML stripping mode setting. Defaults to "index", allowed values are "none", "strip", "index", and "retain".
    AllowEmpty         bool   // Allows empty string to be returned as highlighting result when a snippet could not be generated (no keywords match, or no passages fit the limit). By default, the beginning of original text would be returned instead of an empty string.
    PassageBoundary    string // Ensures that passages do not cross a sentence, paragraph, or zone boundary (when used with an index that has the respective indexing settings enabled). String, allowed values are "sentence", "paragraph", and "zone".
    EmitZones          bool   // Emits an HTML tag with an enclosing zone name before each passage.
}

func (sc *SphinxClient) BuildExcerpts(docs []string, index, words string, opts ExcerptsOpts) (resDocs []string, err error) {
    if len(docs) == 0 {
        return nil, errors.New("BuildExcerpts > Have no documents to process")
    }
    if index == "" {
        return nil, errors.New("BuildExcerpts > index name is empty")
    }
    if words == "" {
        return nil, errors.New("BuildExcerpts > Have no words to highlight")
    }
    if opts.PassageBoundary != "" && opts.PassageBoundary != "sentence" && opts.PassageBoundary != "paragraph" && opts.PassageBoundary != "zone" {
        return nil, fmt.Errorf("BuildExcerpts > PassageBoundary allowed values are 'sentence', 'paragraph', and 'zone', now is: %s", opts.PassageBoundary)
    }

    if opts.BeforeMatch == "" {
        opts.BeforeMatch = "<b>"
    }
    if opts.AfterMatch == "" {
        opts.AfterMatch = "</b>"
    }
    if opts.ChunkSeparator == "" {
        opts.ChunkSeparator = "..."
    }
    if opts.HtmlStripMode == "" {
        opts.HtmlStripMode = "index"
    }
    if opts.Limit == 0 {
        opts.Limit = 256
    }
    if opts.Around == 0 {
        opts.Around = 5
    }
    if opts.StartPassageId == 0 {
        opts.StartPassageId = 1
    }

    var req []byte
    req = writeInt32ToBytes(req, 0)

    iFlags := 1 // remove_spaces
    if opts.ExactPhrase != false {
        iFlags |= 2
    }
    if opts.SinglePassage != false {
        iFlags |= 4
    }
    if opts.UseBoundaries != false {
        iFlags |= 8
    }
    if opts.WeightOrder != false {
        iFlags |= 16
    }
    if opts.QueryMode != false {
        iFlags |= 32
    }
    if opts.ForceAllWords != false {
        iFlags |= 64
    }
    if opts.LoadFiles != false {
        iFlags |= 128
    }
    if opts.AllowEmpty != false {
        iFlags |= 256
    }
    if opts.EmitZones != false {
        iFlags |= 256
    }
    req = writeInt32ToBytes(req, iFlags)
    req = writeLenStrToBytes(req, index)
    req = writeLenStrToBytes(req, words)
    req = writeLenStrToBytes(req, opts.BeforeMatch)
    req = writeLenStrToBytes(req, opts.AfterMatch)
    req = writeLenStrToBytes(req, opts.ChunkSeparator)
    req = writeInt32ToBytes(req, opts.Limit)
    req = writeInt32ToBytes(req, opts.Around)
    req = writeInt32ToBytes(req, opts.LimitPassages)
    req = writeInt32ToBytes(req, opts.LimitWords)
    req = writeInt32ToBytes(req, opts.StartPassageId)
    req = writeLenStrToBytes(req, opts.HtmlStripMode)
    req = writeLenStrToBytes(req, opts.PassageBoundary)
    req = writeInt32ToBytes(req, len(docs))
    for _, doc := range docs {
        req = writeLenStrToBytes(req, doc)
    }
    response, err := sc.doRequest(SearchdCommandExcerpt, SphVerCommandExcerpt, req)
    if err != nil {
        return nil, err
    }

    var bp = byteParser{stream: response}

    resDocs = make([]string, len(docs))
    for i := 0; i < len(docs); i++ {
        resDocs[i] = bp.String()
    }

    return resDocs, nil
}

// UpdateAttributes
// Connect to searchd server and update given attributes on given documents in given indexes.
// values[*][0] is docId, must be an uint64.
// values[*][1:] should be int or []int(mva mode)
// 'ndocs'	-1 on failure, amount of actually found and updated documents (might be 0) on success
func (sc *SphinxClient) UpdateAttributes(index string, attrs []string, values [][]interface{}, ignorenonexistent bool) (ndocs int, err error) {
    if index == "" {
        return -1, errors.New("UpdateAttributes > index name is empty")
    }
    if len(attrs) == 0 {
        return -1, errors.New("UpdateAttributes > no attribute names provided")
    }
    if len(values) < 1 {
        return -1, errors.New("UpdateAttributes > no update values provided")
    }

    for _, v := range values {
        // values[*][0] is docId, so +1
        if len(v) != len(attrs)+1 {
            return -1, fmt.Errorf("UpdateAttributes > update entry has wrong length: %#v", v)
        }
    }

    var mva bool
    if _, ok := values[0][1].([]int); ok {
        mva = true
    }

    // build request
    var req []byte
    req = writeLenStrToBytes(req, index)
    req = writeInt32ToBytes(req, len(attrs))

    if SphVerCommandUpdate > 0x102 {
        if ignorenonexistent {
            req = writeInt32ToBytes(req, 1)
        } else {
            req = writeInt32ToBytes(req, 0)
        }
    }

    for _, attr := range attrs {
        req = writeLenStrToBytes(req, attr)
        if mva {
            req = writeInt32ToBytes(req, 1)
        } else {
            req = writeInt32ToBytes(req, 0)
        }
    }

    req = writeInt32ToBytes(req, len(values))
    for i := 0; i < len(values); i++ {
        if docId, ok := values[i][0].(uint64); !ok {
            return -1, fmt.Errorf("UpdateAttributes > docId must be uint64: %#v", docId)
        } else {
            req = writeInt64ToBytes(req, docId)
        }
        for j := 1; j < len(values[i]); j++ {
            if mva {
                vars, ok := values[i][j].([]int)
                if !ok {
                    return -1, fmt.Errorf("UpdateAttributes > must be []int in mva mode: %#v", vars)
                }
                req = writeInt32ToBytes(req, len(vars))
                for _, v := range vars {
                    req = writeInt32ToBytes(req, v)
                }
            } else {
                v, ok := values[i][j].(int)
                if !ok {
                    return -1, fmt.Errorf("UpdateAttributes > must be int if not in mva mode: %#v", values[i][j])
                }
                req = writeInt32ToBytes(req, v)
            }
        }
    }

    response, er := sc.doRequest(SearchdCommandUpdate, SphVerCommandUpdate, req)
    if er != nil {
        return -1, er
    }

    ndocs = int(binary.BigEndian.Uint32(response[0:4]))
    return
}

type SphinxKeyword struct {
    Tokenized  string
    Normalized string
    Docs       int
    Hits       int
}

func (sc *SphinxClient) BuildKeywords(query, index string, hits bool) (keywords []SphinxKeyword, err error) {
    var req []byte
    req = writeLenStrToBytes(req, query)
    req = writeLenStrToBytes(req, index)
    if hits {
        req = writeInt32ToBytes(req, 1)
    } else {
        req = writeInt32ToBytes(req, 0)
    }
    response, err := sc.doRequest(SearchdCommandKeywords, SphVerCommandKeywords, req)
    if err != nil {
        return nil, err
    }

    var bp = byteParser{stream: response}
    nwords := bp.Int32()
    keywords = make([]SphinxKeyword, nwords)
    for i := 0; i < nwords; i++ {
        var k SphinxKeyword
        k.Tokenized = bp.String()
        k.Normalized = bp.String()
        if hits {
            k.Docs = bp.Int32()
            k.Hits = bp.Int32()
        }
        keywords[i] = k
    }
    return
}

func (sc *SphinxClient) Status() (response [][]string, err error) {
    var req []byte
    req = writeInt32ToBytes(req, 1)

    res, err := sc.doRequest(SearchdCommandStatus, SphVerCommandStatus, req)
    if err != nil {
        return nil, err
    }

    var bp = byteParser{stream: res}
    rows := bp.Uint32()
    cols := bp.Uint32()

    response = make([][]string, rows)
    for i := 0; i < int(rows); i++ {
        response[i] = make([]string, cols)
        for j := 0; j < int(cols); j++ {
            response[i][j] = bp.String()
        }
    }
    return response, nil
}

func (sc *SphinxClient) FlushAttributes() (iFlushTag int, err error) {
    res, err := sc.doRequest(SearchdCommandFlushattrs, SphVerCommandFlushattrs, []byte{})
    if err != nil {
        return -1, err
    }

    if len(res) != 4 {
        return -1, errors.New("FlushAttributes > unexpected response length")
    }

    iFlushTag = int(binary.BigEndian.Uint32(res[0:4]))
    return
}

func (sc *SphinxClient) connect() (err error) {
    if sc.conn != nil {
        return
    }
    sc.connerror = false

    timeout := time.Duration(sc.Timeout) * time.Millisecond

    if sc.Socket != "" {
        if sc.conn, err = net.DialTimeout("unix", sc.Socket, timeout); err != nil {
            sc.connerror = true
            return fmt.Errorf("connect() net.DialTimeout(%d ms) > %v", sc.Timeout, err)
        }
    } else if sc.Port > 0 {
        if sc.conn, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%d", sc.Host, sc.Port), timeout); err != nil {
            sc.connerror = true
            return fmt.Errorf("connect() net.DialTimeout(%d ms) > %v", sc.Timeout, err)
        }
    } else {
        return fmt.Errorf("connect() > No valid socket or port")
    }

    // set deadline
    if err = sc.conn.SetDeadline(time.Now().Add(timeout)); err != nil {
        sc.connerror = true
        return fmt.Errorf("connect() conn.SetDeadline() > %v", err)
    }

    header := make([]byte, 4)
    if _, err = io.ReadFull(sc.conn, header); err != nil {
        sc.connerror = true
        return fmt.Errorf("connect() io.ReadFull() > %v", err)
    }

    version := binary.BigEndian.Uint32(header)
    if version < 1 {
        return fmt.Errorf("connect() > expected searchd protocol version 1+, got version %d", version)
    }

    // send my version
    var i int
    i, err = sc.conn.Write(writeInt32ToBytes([]byte{}, SphVerMajorProto))
    if err != nil {
        sc.connerror = true
        return fmt.Errorf("connect() conn.Write() > %d bytes, %v", i, err)
    }

    // reset deadline
    _ = sc.conn.SetDeadline(time.Time{})

    return
}

func (sc *SphinxClient) Open() (err error) {
    if err = sc.connect(); err != nil {
        return err
    }

    var req []byte
    req = writeInt16ToBytes(req, SearchdCommandPersist)
    req = writeInt16ToBytes(req, 0) // command version
    req = writeInt32ToBytes(req, 4) // body length
    req = writeInt32ToBytes(req, 1) // body

    var n int
    n, err = sc.conn.Write(req)
    if err != nil {
        sc.connerror = true
        return fmt.Errorf("open sc.conn.Write() %d bytes, %s", n, err.Error())
    }

    return nil
}

func (sc *SphinxClient) Close() error {
    if sc.conn == nil {
        return errors.New("client Not connected")
    }

    if err := sc.conn.Close(); err != nil {
        return err
    }

    sc.conn = nil
    return nil
}

func (sc *SphinxClient) doRequest(command int, version int, req []byte) (res []byte, err error) {
    defer func() {
        if x := recover(); x != nil {
            res = nil
            err = fmt.Errorf("doRequest panic > %#v", x)
        }
    }()

    if err = sc.connect(); err != nil {
        return nil, err
    }

    var cmdVerLen []byte
    cmdVerLen = writeInt16ToBytes(cmdVerLen, command)
    cmdVerLen = writeInt16ToBytes(cmdVerLen, version)
    cmdVerLen = writeInt32ToBytes(cmdVerLen, len(req))
    req = append(cmdVerLen, req...)
    _, err = sc.conn.Write(req)
    if err != nil {
        sc.connerror = true
        return nil, fmt.Errorf("conn.Write error: %v", err)
    }

    header := make([]byte, 8)
    if i, err := io.ReadFull(sc.conn, header); err != nil {
        sc.connerror = true
        return nil, fmt.Errorf("doRequest > just read %d bytes into header", i)
    }

    status := binary.BigEndian.Uint16(header[0:2])
    ver := binary.BigEndian.Uint16(header[2:4])
    size := binary.BigEndian.Uint32(header[4:8])
    if size <= 0 {
        return nil, fmt.Errorf("doRequest > invalid response packet size (len=%d)", size)
    }

    res = make([]byte, size)
    if i, err := io.ReadFull(sc.conn, res); err != nil {
        sc.connerror = true
        return nil, fmt.Errorf("doRequest > just read %d bytes into res (size=%d)", i, size)
    }

    switch status {
    case SearchdOk:
    case SearchdWarning:
        wlen := binary.BigEndian.Uint32(res[0:4])
        sc.warning = string(res[4 : 4+wlen])
        res = res[4+wlen:]
    case SearchdError, SearchdRetry:
        wlen := binary.BigEndian.Uint32(res[0:4])
        return nil, fmt.Errorf("doRequest > SEARCHD_ERROR: " + string(res[4:4+wlen]))
    default:
        return nil, fmt.Errorf("doRequest > unknown status code (status=%d), ver: %d", status, ver)
    }

    return res, nil
}

func writeFloat32ToBytes(bs []byte, f float32) []byte {
    buf := new(bytes.Buffer)
    if err := binary.Write(buf, binary.BigEndian, f); err != nil {
        fmt.Println(err)
    }
    return append(bs, buf.Bytes()...)
}

func writeInt16ToBytes(bs []byte, i int) []byte {
    var byte2 = make([]byte, 2)
    binary.BigEndian.PutUint16(byte2, uint16(i))
    return append(bs, byte2...)
}

func writeInt32ToBytes(bs []byte, i int) []byte {
    var byte4 = make([]byte, 4)
    binary.BigEndian.PutUint32(byte4, uint32(i))
    return append(bs, byte4...)
}

func writeInt64ToBytes(bs []byte, ui uint64) []byte {
    var byte8 = make([]byte, 8)
    binary.BigEndian.PutUint64(byte8, ui)
    return append(bs, byte8...)
}

func writeLenStrToBytes(bs []byte, s string) []byte {
    var byte4 = make([]byte, 4)
    binary.BigEndian.PutUint32(byte4, uint32(len(s)))
    bs = append(bs, byte4...)
    return append(bs, []byte(s)...)
}

func EscapeString(s string) string {
    chars := []string{`\`, `(`, `)`, `|`, `-`, `!`, `@`, `~`, `"`, `&`, `/`, `^`, `$`, `=`}
    for _, char := range chars {
        s = strings.Replace(s, char, `\`+char, -1)
    }
    return s
}

func DegreeToRadian(degree float32) float32 {
    return degree * math.Pi / 180
}

// 用于字节数据解析
type byteParser struct {
    stream []byte
    p      int
}

func (bp *byteParser) Int32() (i int) {
    i = int(binary.BigEndian.Uint32(bp.stream[bp.p : bp.p+4]))
    bp.p += 4
    return
}

func (bp *byteParser) Uint32() (i uint32) {
    i = binary.BigEndian.Uint32(bp.stream[bp.p : bp.p+4])
    bp.p += 4
    return
}

func (bp *byteParser) Uint64() (i uint64) {
    i = binary.BigEndian.Uint64(bp.stream[bp.p : bp.p+8])
    bp.p += 8
    return
}

func (bp *byteParser) Float32() (f float32, err error) {
    buf := bytes.NewBuffer(bp.stream[bp.p : bp.p+4])
    bp.p += 4
    if err := binary.Read(buf, binary.BigEndian, &f); err != nil {
        return 0, err
    }
    return f, nil
}

func (bp *byteParser) String() (s string) {
    s = ""
    if slen := bp.Int32(); slen > 0 {
        s = string(bp.stream[bp.p : bp.p+slen])
        bp.p += slen
    }
    return
}
