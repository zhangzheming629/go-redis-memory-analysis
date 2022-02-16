package gorma

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sort"
	"github.com/hhxsv5/go-redis-memory-analysis/storages"
)

type AnalysisConnection struct {
	redis   *storages.RedisClient
	Reports map[string][]Report
}

func NewAnalysisConnection(host string, port uint16, password string) (*AnalysisConnection, error) {
	redis, err := storages.NewRedisClient(host, port, password)
	if err != nil {
		return nil, err
	}
	return &AnalysisConnection{redis, map[string][]Report{}}, nil
}

func (analysis *AnalysisConnection) Close() {
	if analysis.redis != nil {
		_ = analysis.redis.Close()
	}
}

func (analysis AnalysisConnection) Start(limitCount int) {
	fmt.Println("Starting analysis")
	databases, _ := analysis.redis.GetDatabases()
	
	var (
		cursor 					  uint64
		stringReport    		  Report
		listReport      		  Report
		hashReport      		  Report
		setReport       		  Report
		zsetReport      		  Report
		size 					  uint64
		stringSortBySizeReports   SortBySizeReports
		listSortBySizeReports     SortBySizeReports
		hashSortBySizeReports     SortBySizeReports
		setSortBySizeReports     SortBySizeReports
		zsetSortBySizeReports     SortBySizeReports
		totalTypeBySizeReports    SortBySizeReports
		mr     						KeyReports
	)

	analysis.Reports["strType"] = make(SortBySizeReports, 100)
	analysis.Reports["listType"] = make(SortBySizeReports, 100)
	analysis.Reports["hashType"] = make(SortBySizeReports, 100)
	analysis.Reports["setType"] = make(SortBySizeReports, 100)
	analysis.Reports["zsetType"] = make(SortBySizeReports, 100)
	analysis.Reports["totalType"] = make(SortBySizeReports, 100)
    reportMap := make(map[string]SortBySizeReports)

	for db, _ := range databases {
		fmt.Println("Analyzing db", db)
		cursor = 0
		mr = KeyReports{}

		_ = analysis.redis.Select(db)

		for {
			keys, _ := analysis.redis.Scan(&cursor)

			for _, key := range keys {
				size, _ = analysis.redis.GetKeyMemory(key)
				if size < 0 {
                   continue
				}

				keyType, err := analysis.redis.GetKeyType(key)
				if err != nil {
					fmt.Println("analysis getkeyType error", err)
				}
				stringReport = mr[key]
				listReport = mr[key]
				hashReport = mr[key]
				setReport = mr[key]
				zsetReport = mr[key]

				if 0 == strings.Compare("string", keyType) {
					stringReport.Key = key
					stringReport.Type = keyType
					stringReport.Size = size
					stringSortBySizeReports = append(stringSortBySizeReports, stringReport)
					sort.Sort(stringSortBySizeReports)
					stringSortBySizeReports, _ = itertor(limitCount, db, key ,keyType, size, stringReport, stringSortBySizeReports)
					totalTypeBySizeReports = append(totalTypeBySizeReports, stringReport)
				} else if 0 == strings.Compare("list", keyType) {
					listReport.Key = key
					listReport.Type = keyType
					listReport.Size = size
					listSortBySizeReports = append(listSortBySizeReports, listReport)
					sort.Sort(listSortBySizeReports)
					listSortBySizeReports, _ = itertor(limitCount, db, key ,keyType, size, listReport, listSortBySizeReports)
					totalTypeBySizeReports = append(totalTypeBySizeReports, listReport)
				} else if 0 == strings.Compare("hash", keyType) {
					hashReport.Key = key
					hashReport.Type = keyType
					hashReport.Size = size
					hashSortBySizeReports = append(hashSortBySizeReports, hashReport)
					sort.Sort(hashSortBySizeReports)
					hashSortBySizeReports, _ = itertor(limitCount, db, key ,keyType, size, hashReport, hashSortBySizeReports)
					totalTypeBySizeReports = append(totalTypeBySizeReports, hashReport)
				} else if 0 == strings.Compare("set", keyType) {
					setReport.Key = key
					setReport.Type = keyType
					setReport.Size = size
					setSortBySizeReports = append(setSortBySizeReports, setReport)
					sort.Sort(setSortBySizeReports)
					setSortBySizeReports, _ = itertor(limitCount, db, key ,keyType, size, setReport, setSortBySizeReports)
					totalTypeBySizeReports = append(totalTypeBySizeReports, setReport)
				} else {
					zsetReport.Key = key
					zsetReport.Type = keyType
					zsetReport.Size = size
					zsetSortBySizeReports = append(zsetSortBySizeReports, zsetReport)
					sort.Sort(zsetSortBySizeReports)
					zsetSortBySizeReports, _ = itertor(limitCount, db, key ,keyType, size, zsetReport, zsetSortBySizeReports)
					totalTypeBySizeReports = append(totalTypeBySizeReports, zsetReport)
				}
			}

			if cursor == 0 {
				break
			}
		}
	}
	reportMap["strType"] = stringSortBySizeReports
	reportMap["listType"] = listSortBySizeReports
	reportMap["hashType"] = hashSortBySizeReports
	reportMap["setType"] = setSortBySizeReports
	reportMap["zsetType"] = zsetSortBySizeReports
	sort.Sort(totalTypeBySizeReports)

	var hundredLength = 100
    var length = len(totalTypeBySizeReports)
	if len(totalTypeBySizeReports) > 100 {
		var sub = length - hundredLength
		totalTypeBySizeReports = append(totalTypeBySizeReports[:len(totalTypeBySizeReports)-sub], totalTypeBySizeReports[len(totalTypeBySizeReports):]...)
	}
	reportMap["totalType"] = totalTypeBySizeReports

	analysis.Reports["strType"] = reportMap["strType"]
	analysis.Reports["listType"] = reportMap["listType"]
	analysis.Reports["hashType"] = reportMap["hashType"]
	analysis.Reports["setType"] = reportMap["setType"]
	analysis.Reports["zsetType"] = reportMap["zsetType"]
	analysis.Reports["totalType"] = reportMap["totalType"]

}
func itertor(limitCount int, db uint64, key string, keyType string, size uint64, r Report, sr SortBySizeReports) (SortBySizeReports, error) {
  if len(sr) <= limitCount {
	  return sr, nil
  } else if len(sr) > limitCount {
	  sr = append(sr[:len(sr)-1], sr[len(sr):]...)
  }
  return sr,nil
}

func (analysis AnalysisConnection) SaveReports(folder string) error {
	fmt.Println("Saving the results of the analysis into", folder)
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		_ = os.MkdirAll(folder, os.ModePerm)
	}

	var (
		str        string
		jsonReport string
		filename   string
	)
	template := fmt.Sprintf("%s%s%s%s", folder, string(os.PathSeparator), strings.Replace(analysis.redis.Id, ":", "-", -1), ".json")
	filename = template
	fp, err := storages.NewFile(filename, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	if ops,err := json.Marshal(analysis.Reports); err!=nil {
		fmt.Println("json Marshal error: ", err)
	} else {
		jsonReport = string(ops)
	}
	str = fmt.Sprintf("%s\n", jsonReport)
	_, _ = fp.Append([]byte(str))
	fp.Close()
	
	return nil
}
