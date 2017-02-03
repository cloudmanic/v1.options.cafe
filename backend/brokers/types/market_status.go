package types

type MarketStatus struct {
  Date string 
  State string
  Description string      
}

/*
  <date>2014-05-27</date>
  <description>Market is open from 09:30 to 16:00.</description>
  <next_change>16:00</next_change>
  <next_state>postmarket</next_state>
  <state>open</state>
  <timestamp>1401219046</timestamp>
*/