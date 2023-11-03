package user

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"smatyx.com/shared/cast"
	"smatyx.com/shared/database"
	"smatyx.com/shared/entities"
	"smatyx.com/shared/meta"
	"smatyx.com/shared/server"
)

type UserSearchReq struct {
	Id     []uuid.UUID
	Limit  int
	Offset int
}

type UserSearchRow struct {
	Id          string             `json:"id"`
	Email       string             `json:"email"`
	Phone       string             `json:"phone"`
	FirstName   string             `json:"firstName"`
	LastName    string             `json:"lastName"`
	State       entities.UserState `json:"state"`
	CreatedTime time.Time          `json:"createdTime"`
	UpdatedTime time.Time          `json:"updatedTime"`
}

func UserSearch(txt *server.Transaction) {
	urlQuery := txt.HttpRequest.URL.Query()

	req := UserSearchReq{
		Id:     make([]uuid.UUID, 0, 10),
		Limit:  20,
		Offset: 0,
	}

	idStrings := urlQuery["id"]
	fmt.Print(idStrings)
	for _, idString := range idStrings {
		id, err := uuid.Parse(idString)
		if err != nil {
			panic(err)
		}
		req.Id = append(req.Id, id)
	}

	fmt.Print(req.Id)
	selectQb := database.NewSelectBuilder([]string{
		meta.User_Id,
		meta.User_Email,
		meta.User_Phone,
		meta.User_FirstName,
		meta.User_LastName,
		meta.User_State,
		meta.User_CreatedTime,
		meta.User_UpdatedTime,
	})
	countQb := database.NewSelectBuilder([]string{"count(*)"})

	qbClausesBuild := func(qb *database.QueryBuilder) {
		qb.From(meta.User.Name)
		qb.WhereInVars(meta.User.Name, meta.User.FieldNames[meta.User_Id_Idx], cast.CopyToSliceAny(req.Id))
		qb.LimitVar(req.Limit)
		qb.OffsetVar(req.Offset)
	}

	qbClausesBuild(selectQb)
	qbClausesBuild(countQb)

	users := make([]UserSearchRow, 0, req.Limit)
	count := 0

	userRowProcess := func(pgRep *database.PgRep) error {
		row := UserSearchRow{}
		
		err := pgRep.Scan(
			&row.Id,
			&row.Email,
			&row.Phone,
			&row.FirstName,
			&row.LastName,
			&row.State,
			&row.CreatedTime,
			&row.UpdatedTime)
		if err != nil {
			return err
		}

		users = append(users, row)
		fmt.Print(1, users)
		return nil
	}

	countProcess := func(pgRep *database.PgRep) error {
		err := pgRep.Scan(&count)
		if err != nil {
			return err
		}
		return nil
	}

	queries := database.NewQueriesFromBuilders(
		txt.Context,
		[]*database.QueryBuilder{selectQb, countQb},
		[]database.ProcessRepFunc{userRowProcess},
	)
	queries.Submit()

	txt.WriteBytes(server.PageRepJson(count, req.Offset, req.Limit, users, nil))
}
