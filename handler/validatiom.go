package handler 

func isNotFoundError(err error)bool{

	switch err.Error(){
	case "sql: no rows in result set": 
	case "user ID does not exist":
	default:
		return  false 
	}
	return true

}