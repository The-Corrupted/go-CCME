/*

CCMESQL server. Used to store the results of tests and serve them back up to you if
trying to view the results page. Uses a basic datamap that is used in server to 
create the results table.

*/

package CCMESql


//TODO: Figure out a way to setup a database if one does not presently exist
/* Current Setup: Create psql database by the name of analysisquota. 
   Create table by the name of results that has id int, videoName varchar, badFrames int, extraFrames int,
   frameNumber int and finalResult varchar
   */
   
import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"github.com/The-Corrupted/go-CCME/Helpers/OsHandler"
)

type DataMap struct {
	Id string
	VideoName string
	BadFrames string
	ExtraFrames string
	FrameNumber string
	FinalResult string
	Fulloutput string
}

//SQL Login Info: TODO Find a more secure way to store sql login information.

// For Simple video analysis

func GetCredentials() string {
	creds, err := ioutil.ReadFile(fmt.Sprintf("%s/.creds", OsHandler.SetUserDir()))
	if err != nil {
		fmt.Println(err)
		return "Failed"
	}
	return string(creds)
}

func ReturnCredentials() string {
	
	return GetCredentials()
}

func Retrieve_Data() []DataMap {
	FullData := make([]DataMap, 0)
	var loginInfo = ReturnCredentials()
	db, err := sql.Open("postgres", loginInfo)
	if err != nil {
		fmt.Println("Failed To Open DB")
	} else {
		rows, err := db.Query("SELECT id, videoName, badFrames, extraFrames, frameNumber, finalResult FROM results ORDER BY id DESC LIMIT 1000")
		if err != nil {
			fmt.Printf("Failed to query database: %v", err)

		} else {
			defer rows.Close()
			for rows.Next() {
				var dat DataMap
				rows.Scan(&dat.Id, &dat.VideoName, &dat.BadFrames, &dat.ExtraFrames, &dat.FrameNumber,
					     &dat.FinalResult)
				FullData = append(FullData, dat)

			}
		}
		return FullData
	}
	FailData := make([]DataMap,0)
	return FailData
}

func Last_Id() int64 {
	var lastID int64 = 0
	var loginInfo = ReturnCredentials()
	db, err := sql.Open("postgres", loginInfo)
	if err != nil {
		fmt.Printf("Failed to open DB: %v\n", err)
	} else {
		rows, err := db.Query(` SELECT id FROM results ORDER BY id DESC limit 1`)
		if err != nil {
			fmt.Printf("Failed to query database: %v\n", err)
		} else {
			defer rows.Close()
			for rows.Next() {
				rows.Scan(&lastID)
			}
		}
		return lastID
	}
	return lastID
}

func Retrieve_Fulloutput(id string) string {
	var Fulloutput string = ""
	var loginInfo string = ReturnCredentials()
	db, err := sql.Open("postgres", loginInfo)
	if ( err != nil ) {
		fmt.Println("Failed To Open DB")
	} else {
		rows, err := db.Query(fmt.Sprintf("SELECT fulloutput FROM results WHERE id = %s", id))
		if ( err != nil ) {
			fmt.Printf("Failed to query database: %v", err)

		} else {
			defer rows.Close()
			for rows.Next() {
				rows.Scan(&Fulloutput)
			}
			fmt.Printf("Retrieved fulloutput: %s\n", Fulloutput)
			}
		}
	return Fulloutput;
}

func UpdateDB(id int64, videoName string, badFrames uint64, extraFrames uint64, frameNumber uint64, finalResult string, fulloutput string) (bool, error) {
	var loginInfo = ReturnCredentials()
	db, err := sql.Open("postgres", loginInfo)
	if err != nil {
		return false, err
	}
	sqlStatement := `INSERT INTO results(id, videoName, badFrames, extraFrames, frameNumber, finalResult, fulloutput)
						 VALUES ($1, $2, $3, $4, $5, $6, $7);`
	_, err = db.Exec(sqlStatement, id, videoName, badFrames, extraFrames, frameNumber, finalResult, fulloutput)
	if err != nil {
		return false, err
	}
	fmt.Println("DB updated successfully.")
	return true, nil
}