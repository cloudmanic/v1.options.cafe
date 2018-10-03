package types

type UserProfile struct {
	Id       string
	Name     string
	Accounts []UserProfileAccounts
}

type UserProfileAccounts struct {
	AccountNumber  string
	Classification string
	DayTrader      bool
	OptionLevel    int
	Status         string
	Type           string
}

/*
<profile>
  <id>id-gcostanza</id>
  <account>
    <account_number>12345678</account_number>
    <classification>individual</classification>
    <day_trader>true</day_trader>
    <option_level>5</option_level>
    <status>Approved</status>
    <type>pdt</type>
    <last_update_date>2014-06-23T15:49:05.847Z</last_update_date>
  </account>
  <account>
    <account_number>87654321</account_number>
    <classification>individual</classification>
    <day_trader>false</day_trader>
    <option_level>3</option_level>
    <status>Approved</status>
    <type>margin</type>
    <last_update_date>2014-06-23T15:49:05.868Z</last_update_date>
  </account>
  <name>George Costanza</name>
</profile>
*/
