//
// Date: 11/7/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package tradier

import (
	"testing"

	"github.com/nbio/st"
	gock "gopkg.in/h2non/gock.v1"
)

//
// Search by company name or symbol
//
func TestSearchBySymbolOrCompanyName(t *testing.T) {

	// Flush pending mocks after test execution
	defer gock.Off()

	// Setup mock request.
	gock.New(apiBaseUrl).
		Get("/markets/lookup").
		Reply(200).
		BodyString(`{"securities":{"security":[{"symbol":"SBUX","exchange":"Q","type":"stock","description":"Starbucks Corp"},{"symbol":"SBUM","exchange":"V","type":"stock","description":"Silver Buckle Mines Inc"},{"symbol":"SBUS","exchange":"P","type":"etf","description":"ETFS Diversified-Factor U.S. Large Cap Index Fund"},{"symbol":"SBUY","exchange":"V","type":"stock","description":null}]}}`)

	// Setup mock request.
	gock.New(apiBaseUrl).
		Get("/markets/search").
		Reply(200).
		BodyString(`{"securities":{"security":[{"symbol":"SBUX","exchange":"Q","type":"stock","description":"Starbucks Corp"},{"symbol":"STSC","exchange":"V","type":"stock","description":"Start Scientific Inc"},{"symbol":"STRZA","exchange":"Q","type":"stock","description":"Starz"},{"symbol":"SFR","exchange":"N","type":"stock","description":"Starwood Waypoint Homes"},{"symbol":"STWD","exchange":"N","type":"stock","description":"Starwood Property Trust Inc"},{"symbol":"SBLK","exchange":"Q","type":"stock","description":"Star Bulk Carriers Corp."},{"symbol":"SRDP","exchange":"V","type":"stock","description":"Star Resorts Development Inc"},{"symbol":"SRT","exchange":"N","type":"stock","description":"StarTek Inc"},{"symbol":"SGU","exchange":"N","type":"stock","description":"Star Group LP"},{"symbol":"STAU","exchange":"V","type":"stock","description":"Star Nutrition Inc"},{"symbol":"SALY","exchange":"V","type":"stock","description":"Star Ally Inc"},{"symbol":"STGZ","exchange":"V","type":"stock","description":"Stargaze Entertainment Group Inc"},{"symbol":"STRZ","exchange":"V","type":"stock","description":"Star Buffet Inc"},{"symbol":"STRH","exchange":"V","type":"stock","description":"Star8 Corp"},{"symbol":"SMRS","exchange":"V","type":"stock","description":"Star Mountain Resources Inc"},{"symbol":"SRGZ","exchange":"V","type":"stock","description":"Star Gold Corp."},{"symbol":"SSET","exchange":"V","type":"stock","description":"Starstream Entertainment Inc"},{"symbol":"SCPD","exchange":"V","type":"stock","description":"Star Century Pandaho Corp"},{"symbol":"SRSK","exchange":"V","type":"stock","description":"Star Struck Ltd"},{"symbol":"SFIGA","exchange":"V","type":"stock","description":"Star Financial Group Inc"},{"symbol":"STNX","exchange":"V","type":"stock","description":"Startronix International Inc"},{"symbol":"SWMD","exchange":"V","type":"stock","description":"Starwin Media Holdings Inc"},{"symbol":"SETY","exchange":"V","type":"stock","description":"Star Entertainment Group Inc"},{"symbol":"SREN","exchange":"V","type":"stock","description":"Star Entertainment Group"},{"symbol":"SRHBF","exchange":"V","type":"stock","description":"STARHUB LTD ORD by Starhub, Ltd."},{"symbol":"EHGRF","exchange":"V","type":"stock","description":"STAR ENTMT GROUP LTD ORD by STAR ENTMT GROUP LTD"},{"symbol":"SLSC","exchange":"V","type":"stock","description":"Starlight Supply Chain Management Company"},{"symbol":"SRVTF","exchange":"V","type":"stock","description":"STAR VAULT AB MALMO (B) by Star Vault AB, Malmo"},{"symbol":"SHVLD","exchange":"V","type":"stock","description":"Starcore International Mines Ltd"},{"symbol":"SPHRF","exchange":"V","type":"stock","description":"STARPHARMA HOLDINGS LTD by Starpharma Holdings Ltd."},{"symbol":"SWOOF","exchange":"V","type":"stock","description":"STARWOOD EUROPEAN RE FIN by StarWood European Real Estate Fin Ltd"},{"symbol":"LSGEF","exchange":"V","type":"stock","description":"STARR PEAK EXPL LTD by STARR PEAK EXPL LTD"},{"symbol":"SRHBY","exchange":"V","type":"stock","description":"STARHUB LTD UNSP\/ADR by Starhub, Ltd."},{"symbol":"STBEF","exchange":"V","type":"stock","description":"STARBREEZE AB ORD by STARBREEZE AB"},{"symbol":"SPHRY","exchange":"V","type":"stock","description":"STARPHARMA HLDGS S\/ADR by Starpharma Holdings Ltd."},{"symbol":"STIV","exchange":"V","type":"stock","description":"STARINVEST GROUP INC by StarInvest Group, Inc."},{"symbol":"SAEC","exchange":"V","type":"stock","description":"STARLIGHT ENERGY CORP by Starlight Energy Corp."},{"symbol":"SRFDF","exchange":"V","type":"stock","description":"STARFIELD RES INC by Starfield Resources Inc."},{"symbol":"STRZB","exchange":"Q","type":"stock","description":"STARZ COM SER B"},{"symbol":"STXMF","exchange":"V","type":"stock","description":"STARREX INTERNATIONAL LTD by Starrex International Ltd."},{"symbol":"SATLF","exchange":"V","type":"stock","description":"START TODAY CO LTD ORD by Start Today Co., Ltd."},{"symbol":"STFK","exchange":"U","type":"stock","description":"STARFLICK.COM by Starflick.com"},{"symbol":"SNAVF","exchange":"V","type":"stock","description":"STAR NAVIGATION SY GP LTD by Star Navigation Systems Group Ltd."},{"symbol":"SBLKL","exchange":"Q","type":"stock","description":"Star Bulk Carriers Corp. - 8.00% Senior Notes Due 2019"},{"symbol":"SCXLB","exchange":"V","type":"stock","description":"STARRETT (L S) CO B COM by Starrett (L.S.) Co. (The)"},{"symbol":"SRBGF","exchange":"V","type":"stock","description":"STARRAG RORSCHACHERBERG N by Starrag Rorschacherberg"},{"symbol":"SKKB","exchange":"V","type":"stock","description":"Stark Naked Bobbers"},{"symbol":"SGLMF","exchange":"V","type":"stock","description":"STARHILL GLOBAL RE INV TR by Starhill Global Real Estate Investment Trust"},{"symbol":"SHVLF","exchange":"V","type":"stock","description":"Starcore International Mines Ltd"},{"symbol":"106677:OOTC","exchange":"V","type":"stock","description":"STARINVEST GROUP INC by StarInvest Group, Inc."},{"symbol":"SLEIY","exchange":"V","type":"stock","description":"STARLIGHT INTL HLDGS ADR by Starlight International Holdings Ltd."},{"symbol":"SRTTY","exchange":"V","type":"stock","description":"START TODAY LTD UNSP\/ADR by Start Today Co., Ltd."},{"symbol":"SGGGF","exchange":"V","type":"stock","description":"STARGROUP LTD ORD by Stargroup Ltd"},{"symbol":"STMDF","exchange":"V","type":"stock","description":"STARTMONDAY TECHNOLOGY CO by StartMonday Technology Corp."},{"symbol":"SAKYF","exchange":"V","type":"stock","description":"STAR MICA CO LTD ORD by Star Mica Co., Ltd."},{"symbol":"ATRKD","exchange":"V","type":"stock","description":"STAR ALLIANCE INTL CORP by Star Alliance International Corp."},{"symbol":"KMTGD","exchange":"V","type":"stock","description":"StarPower ON Systems Inc"},{"symbol":"STFKE","exchange":"V","type":"stock","description":"Starflick.com"},{"symbol":"LSGED","exchange":"V","type":"stock","description":"Starr Peak Exploration Ltd"},{"symbol":"HOT","exchange":"N","type":"stock","description":"Starwood Hotels & Resorts Worldwide, Inc. Common Stock"},{"symbol":"STAL","exchange":"V","type":"stock","description":"STAR ALLIANCE INTL CORP by Star Alliance International Corp."},{"symbol":"SPOS","exchange":"V","type":"stock","description":"STARPOWER ON SYSTEMS INC by StarPower ON Systems, Inc."},{"symbol":"STCB","exchange":"V","type":"stock","description":"STARCO BRANDS INC by Starco Brands, Inc."}]}}
  [{SBUX Starbucks Corp} {STSC Start Scientific Inc} {STRZA Starz} {SFR Starwood Waypoint Homes} {STWD Starwood Property Trust Inc} {SBLK Star Bulk Carriers Corp.} {SRDP Star Resorts Development Inc} {SRT StarTek Inc} {SGU Star Group LP} {STAU Star Nutrition Inc} {SALY Star Ally Inc} {STGZ Stargaze Entertainment Group Inc} {STRZ Star Buffet Inc} {STRH Star8 Corp} {SMRS Star Mountain Resources Inc} {SRGZ Star Gold Corp.} {SSET Starstream Entertainment Inc} {SCPD Star Century Pandaho Corp} {SRSK Star Struck Ltd} {SFIGA Star Financial Group Inc} {STNX Startronix International Inc} {SWMD Starwin Media Holdings Inc} {SETY Star Entertainment Group Inc} {SREN Star Entertainment Group} {SRHBF STARHUB LTD ORD by Starhub, Ltd.} {EHGRF STAR ENTMT GROUP LTD ORD by STAR ENTMT GROUP LTD} {SLSC Starlight Supply Chain Management Company} {SRVTF STAR VAULT AB MALMO (B) by Star Vault AB, Malmo} {SHVLD Starcore International Mines Ltd} {SPHRF STARPHARMA HOLDINGS LTD by Starpharma Holdings Ltd.} {SWOOF STARWOOD EUROPEAN RE FIN by StarWood European Real Estate Fin Ltd} {LSGEF STARR PEAK EXPL LTD by STARR PEAK EXPL LTD} {SRHBY STARHUB LTD UNSP/ADR by Starhub, Ltd.} {STBEF STARBREEZE AB ORD by STARBREEZE AB} {SPHRY STARPHARMA HLDGS S/ADR by Starpharma Holdings Ltd.} {STIV STARINVEST GROUP INC by StarInvest Group, Inc.} {SAEC STARLIGHT ENERGY CORP by Starlight Energy Corp.} {SRFDF STARFIELD RES INC by Starfield Resources Inc.} {STRZB STARZ COM SER B} {STXMF STARREX INTERNATIONAL LTD by Starrex International Ltd.} {SATLF START TODAY CO LTD ORD by Start Today Co., Ltd.} {STFK STARFLICK.COM by Starflick.com} {SNAVF STAR NAVIGATION SY GP LTD by Star Navigation Systems Group Ltd.} {SBLKL Star Bulk Carriers Corp. - 8.00% Senior Notes Due 2019} {SCXLB STARRETT (L S) CO B COM by Starrett (L.S.) Co. (The)} {SRBGF STARRAG RORSCHACHERBERG N by Starrag Rorschacherberg} {SKKB Stark Naked Bobbers} {SGLMF STARHILL GLOBAL RE INV TR by Starhill Global Real Estate Investment Trust} {SHVLF Starcore International Mines Ltd} {106677:OOTC STARINVEST GROUP INC by StarInvest Group, Inc.} {SLEIY STARLIGHT INTL HLDGS ADR by Starlight International Holdings Ltd.} {SRTTY START TODAY LTD UNSP/ADR by Start Today Co., Ltd.} {SGGGF STARGROUP LTD ORD by Stargroup Ltd} {STMDF STARTMONDAY TECHNOLOGY CO by StartMonday Technology Corp.} {SAKYF STAR MICA CO LTD ORD by Star Mica Co., Ltd.} {ATRKD STAR ALLIANCE INTL CORP by Star Alliance International Corp.} {KMTGD StarPower ON Systems Inc} {STFKE Starflick.com} {LSGED Starr Peak Exploration Ltd} {HOT Starwood Hotels & Resorts Worldwide, Inc. Common Stock} {STAL STAR ALLIANCE INTL CORP by Star Alliance International Corp.} {SPOS STARPOWER ON SYSTEMS INC by StarPower ON Systems, Inc.} {STCB STARCO BRANDS INC by Starco Brands, Inc.}]`)

	// Create new tradier isntance
	tradier := &Api{}

	// Make API call
	symbols, err := tradier.SearchBySymbolOrCompanyName("star")

	if err != nil {
		panic(err)
	}

	// Verify the data was return as expected
	st.Expect(t, len(symbols), 66)
	st.Expect(t, symbols[0].Name, "SBUX")
	st.Expect(t, symbols[0].Description, "Starbucks Corp")
	st.Expect(t, symbols[1].Name, "SBUM")
	st.Expect(t, symbols[1].Description, "Silver Buckle Mines Inc")

	// Verify that we don't have pending mocks
	st.Expect(t, gock.IsDone(), true)

}

//
// Search by company name
//
func TestSearchBySymbolName(t *testing.T) {

	// Flush pending mocks after test execution
	defer gock.Off()

	// Setup mock request.
	gock.New(apiBaseUrl).
		Get("/markets/lookup").
		Reply(200).
		BodyString(`{"securities":{"security":[{"symbol":"SBUX","exchange":"Q","type":"stock","description":"Starbucks Corp"},{"symbol":"SBUM","exchange":"V","type":"stock","description":"Silver Buckle Mines Inc"},{"symbol":"SBUS","exchange":"P","type":"etf","description":"ETFS Diversified-Factor U.S. Large Cap Index Fund"},{"symbol":"SBUY","exchange":"V","type":"stock","description":null}]}}`)

		// Create new tradier isntance
	tradier := &Api{}

	// Make API call
	symbols, err := tradier.SearchBySymbolName("sbu")

	if err != nil {
		panic(err)
	}

	// Verify the data was return as expected
	st.Expect(t, len(symbols), 4)
	st.Expect(t, symbols[0].Name, "SBUX")
	st.Expect(t, symbols[0].Description, "Starbucks Corp")
	st.Expect(t, symbols[1].Name, "SBUM")
	st.Expect(t, symbols[1].Description, "Silver Buckle Mines Inc")

	// Verify that we don't have pending mocks
	st.Expect(t, gock.IsDone(), true)

	// ----------- Test One Result -------------- //

	// Setup mock request.
	gock.New(apiBaseUrl).
		Get("/markets/lookup").
		Reply(200).
		BodyString(`{"securities":{"security":{"symbol":"SBUX","exchange":"Q","type":"stock","description":"Starbucks Corp"}}}`)

	// Make API call
	symbols2, err := tradier.SearchBySymbolName("sbux")

	if err != nil {
		panic(err)
	}

	// Verify the data was return as expected
	st.Expect(t, len(symbols2), 1)
	st.Expect(t, symbols2[0].Name, "SBUX")
	st.Expect(t, symbols2[0].Description, "Starbucks Corp")

	// Verify that we don't have pending mocks
	st.Expect(t, gock.IsDone(), true)
}

//
// Search by company name
//
func TestSearchByCompanyName(t *testing.T) {

	// Flush pending mocks after test execution
	defer gock.Off()

	// Setup mock request.
	gock.New(apiBaseUrl).
		Get("/markets/search").
		Reply(200).
		BodyString(`{"securities":{"security":[{"symbol":"SBUX","exchange":"Q","type":"stock","description":"Starbucks Corp"},{"symbol":"STSC","exchange":"V","type":"stock","description":"Start Scientific Inc"},{"symbol":"STRZA","exchange":"Q","type":"stock","description":"Starz"},{"symbol":"SFR","exchange":"N","type":"stock","description":"Starwood Waypoint Homes"},{"symbol":"STWD","exchange":"N","type":"stock","description":"Starwood Property Trust Inc"},{"symbol":"SBLK","exchange":"Q","type":"stock","description":"Star Bulk Carriers Corp."},{"symbol":"SRDP","exchange":"V","type":"stock","description":"Star Resorts Development Inc"},{"symbol":"SRT","exchange":"N","type":"stock","description":"StarTek Inc"},{"symbol":"SGU","exchange":"N","type":"stock","description":"Star Group LP"},{"symbol":"STAU","exchange":"V","type":"stock","description":"Star Nutrition Inc"},{"symbol":"SALY","exchange":"V","type":"stock","description":"Star Ally Inc"},{"symbol":"STGZ","exchange":"V","type":"stock","description":"Stargaze Entertainment Group Inc"},{"symbol":"STRZ","exchange":"V","type":"stock","description":"Star Buffet Inc"},{"symbol":"STRH","exchange":"V","type":"stock","description":"Star8 Corp"},{"symbol":"SMRS","exchange":"V","type":"stock","description":"Star Mountain Resources Inc"},{"symbol":"SRGZ","exchange":"V","type":"stock","description":"Star Gold Corp."},{"symbol":"SSET","exchange":"V","type":"stock","description":"Starstream Entertainment Inc"},{"symbol":"SCPD","exchange":"V","type":"stock","description":"Star Century Pandaho Corp"},{"symbol":"SRSK","exchange":"V","type":"stock","description":"Star Struck Ltd"},{"symbol":"SFIGA","exchange":"V","type":"stock","description":"Star Financial Group Inc"},{"symbol":"STNX","exchange":"V","type":"stock","description":"Startronix International Inc"},{"symbol":"SWMD","exchange":"V","type":"stock","description":"Starwin Media Holdings Inc"},{"symbol":"SETY","exchange":"V","type":"stock","description":"Star Entertainment Group Inc"},{"symbol":"SREN","exchange":"V","type":"stock","description":"Star Entertainment Group"},{"symbol":"SRHBF","exchange":"V","type":"stock","description":"STARHUB LTD ORD by Starhub, Ltd."},{"symbol":"EHGRF","exchange":"V","type":"stock","description":"STAR ENTMT GROUP LTD ORD by STAR ENTMT GROUP LTD"},{"symbol":"SLSC","exchange":"V","type":"stock","description":"Starlight Supply Chain Management Company"},{"symbol":"SRVTF","exchange":"V","type":"stock","description":"STAR VAULT AB MALMO (B) by Star Vault AB, Malmo"},{"symbol":"SHVLD","exchange":"V","type":"stock","description":"Starcore International Mines Ltd"},{"symbol":"SPHRF","exchange":"V","type":"stock","description":"STARPHARMA HOLDINGS LTD by Starpharma Holdings Ltd."},{"symbol":"SWOOF","exchange":"V","type":"stock","description":"STARWOOD EUROPEAN RE FIN by StarWood European Real Estate Fin Ltd"},{"symbol":"LSGEF","exchange":"V","type":"stock","description":"STARR PEAK EXPL LTD by STARR PEAK EXPL LTD"},{"symbol":"SRHBY","exchange":"V","type":"stock","description":"STARHUB LTD UNSP\/ADR by Starhub, Ltd."},{"symbol":"STBEF","exchange":"V","type":"stock","description":"STARBREEZE AB ORD by STARBREEZE AB"},{"symbol":"SPHRY","exchange":"V","type":"stock","description":"STARPHARMA HLDGS S\/ADR by Starpharma Holdings Ltd."},{"symbol":"STIV","exchange":"V","type":"stock","description":"STARINVEST GROUP INC by StarInvest Group, Inc."},{"symbol":"SAEC","exchange":"V","type":"stock","description":"STARLIGHT ENERGY CORP by Starlight Energy Corp."},{"symbol":"SRFDF","exchange":"V","type":"stock","description":"STARFIELD RES INC by Starfield Resources Inc."},{"symbol":"STRZB","exchange":"Q","type":"stock","description":"STARZ COM SER B"},{"symbol":"STXMF","exchange":"V","type":"stock","description":"STARREX INTERNATIONAL LTD by Starrex International Ltd."},{"symbol":"SATLF","exchange":"V","type":"stock","description":"START TODAY CO LTD ORD by Start Today Co., Ltd."},{"symbol":"STFK","exchange":"U","type":"stock","description":"STARFLICK.COM by Starflick.com"},{"symbol":"SNAVF","exchange":"V","type":"stock","description":"STAR NAVIGATION SY GP LTD by Star Navigation Systems Group Ltd."},{"symbol":"SBLKL","exchange":"Q","type":"stock","description":"Star Bulk Carriers Corp. - 8.00% Senior Notes Due 2019"},{"symbol":"SCXLB","exchange":"V","type":"stock","description":"STARRETT (L S) CO B COM by Starrett (L.S.) Co. (The)"},{"symbol":"SRBGF","exchange":"V","type":"stock","description":"STARRAG RORSCHACHERBERG N by Starrag Rorschacherberg"},{"symbol":"SKKB","exchange":"V","type":"stock","description":"Stark Naked Bobbers"},{"symbol":"SGLMF","exchange":"V","type":"stock","description":"STARHILL GLOBAL RE INV TR by Starhill Global Real Estate Investment Trust"},{"symbol":"SHVLF","exchange":"V","type":"stock","description":"Starcore International Mines Ltd"},{"symbol":"106677:OOTC","exchange":"V","type":"stock","description":"STARINVEST GROUP INC by StarInvest Group, Inc."},{"symbol":"SLEIY","exchange":"V","type":"stock","description":"STARLIGHT INTL HLDGS ADR by Starlight International Holdings Ltd."},{"symbol":"SRTTY","exchange":"V","type":"stock","description":"START TODAY LTD UNSP\/ADR by Start Today Co., Ltd."},{"symbol":"SGGGF","exchange":"V","type":"stock","description":"STARGROUP LTD ORD by Stargroup Ltd"},{"symbol":"STMDF","exchange":"V","type":"stock","description":"STARTMONDAY TECHNOLOGY CO by StartMonday Technology Corp."},{"symbol":"SAKYF","exchange":"V","type":"stock","description":"STAR MICA CO LTD ORD by Star Mica Co., Ltd."},{"symbol":"ATRKD","exchange":"V","type":"stock","description":"STAR ALLIANCE INTL CORP by Star Alliance International Corp."},{"symbol":"KMTGD","exchange":"V","type":"stock","description":"StarPower ON Systems Inc"},{"symbol":"STFKE","exchange":"V","type":"stock","description":"Starflick.com"},{"symbol":"LSGED","exchange":"V","type":"stock","description":"Starr Peak Exploration Ltd"},{"symbol":"HOT","exchange":"N","type":"stock","description":"Starwood Hotels & Resorts Worldwide, Inc. Common Stock"},{"symbol":"STAL","exchange":"V","type":"stock","description":"STAR ALLIANCE INTL CORP by Star Alliance International Corp."},{"symbol":"SPOS","exchange":"V","type":"stock","description":"STARPOWER ON SYSTEMS INC by StarPower ON Systems, Inc."},{"symbol":"STCB","exchange":"V","type":"stock","description":"STARCO BRANDS INC by Starco Brands, Inc."}]}}
	[{SBUX Starbucks Corp} {STSC Start Scientific Inc} {STRZA Starz} {SFR Starwood Waypoint Homes} {STWD Starwood Property Trust Inc} {SBLK Star Bulk Carriers Corp.} {SRDP Star Resorts Development Inc} {SRT StarTek Inc} {SGU Star Group LP} {STAU Star Nutrition Inc} {SALY Star Ally Inc} {STGZ Stargaze Entertainment Group Inc} {STRZ Star Buffet Inc} {STRH Star8 Corp} {SMRS Star Mountain Resources Inc} {SRGZ Star Gold Corp.} {SSET Starstream Entertainment Inc} {SCPD Star Century Pandaho Corp} {SRSK Star Struck Ltd} {SFIGA Star Financial Group Inc} {STNX Startronix International Inc} {SWMD Starwin Media Holdings Inc} {SETY Star Entertainment Group Inc} {SREN Star Entertainment Group} {SRHBF STARHUB LTD ORD by Starhub, Ltd.} {EHGRF STAR ENTMT GROUP LTD ORD by STAR ENTMT GROUP LTD} {SLSC Starlight Supply Chain Management Company} {SRVTF STAR VAULT AB MALMO (B) by Star Vault AB, Malmo} {SHVLD Starcore International Mines Ltd} {SPHRF STARPHARMA HOLDINGS LTD by Starpharma Holdings Ltd.} {SWOOF STARWOOD EUROPEAN RE FIN by StarWood European Real Estate Fin Ltd} {LSGEF STARR PEAK EXPL LTD by STARR PEAK EXPL LTD} {SRHBY STARHUB LTD UNSP/ADR by Starhub, Ltd.} {STBEF STARBREEZE AB ORD by STARBREEZE AB} {SPHRY STARPHARMA HLDGS S/ADR by Starpharma Holdings Ltd.} {STIV STARINVEST GROUP INC by StarInvest Group, Inc.} {SAEC STARLIGHT ENERGY CORP by Starlight Energy Corp.} {SRFDF STARFIELD RES INC by Starfield Resources Inc.} {STRZB STARZ COM SER B} {STXMF STARREX INTERNATIONAL LTD by Starrex International Ltd.} {SATLF START TODAY CO LTD ORD by Start Today Co., Ltd.} {STFK STARFLICK.COM by Starflick.com} {SNAVF STAR NAVIGATION SY GP LTD by Star Navigation Systems Group Ltd.} {SBLKL Star Bulk Carriers Corp. - 8.00% Senior Notes Due 2019} {SCXLB STARRETT (L S) CO B COM by Starrett (L.S.) Co. (The)} {SRBGF STARRAG RORSCHACHERBERG N by Starrag Rorschacherberg} {SKKB Stark Naked Bobbers} {SGLMF STARHILL GLOBAL RE INV TR by Starhill Global Real Estate Investment Trust} {SHVLF Starcore International Mines Ltd} {106677:OOTC STARINVEST GROUP INC by StarInvest Group, Inc.} {SLEIY STARLIGHT INTL HLDGS ADR by Starlight International Holdings Ltd.} {SRTTY START TODAY LTD UNSP/ADR by Start Today Co., Ltd.} {SGGGF STARGROUP LTD ORD by Stargroup Ltd} {STMDF STARTMONDAY TECHNOLOGY CO by StartMonday Technology Corp.} {SAKYF STAR MICA CO LTD ORD by Star Mica Co., Ltd.} {ATRKD STAR ALLIANCE INTL CORP by Star Alliance International Corp.} {KMTGD StarPower ON Systems Inc} {STFKE Starflick.com} {LSGED Starr Peak Exploration Ltd} {HOT Starwood Hotels & Resorts Worldwide, Inc. Common Stock} {STAL STAR ALLIANCE INTL CORP by Star Alliance International Corp.} {SPOS STARPOWER ON SYSTEMS INC by StarPower ON Systems, Inc.} {STCB STARCO BRANDS INC by Starco Brands, Inc.}]`)

		// Create new tradier instance
	tradier := &Api{}

	// Make API call
	symbols, err := tradier.SearchByCompanyName("star")

	if err != nil {
		panic(err)
	}

	// Verify the data was return as expected
	st.Expect(t, len(symbols), 63)
	st.Expect(t, symbols[0].Name, "SBUX")
	st.Expect(t, symbols[0].Description, "Starbucks Corp")
	st.Expect(t, symbols[1].Name, "STSC")
	st.Expect(t, symbols[1].Description, "Start Scientific Inc")

	// Verify that we don't have pending mocks
	st.Expect(t, gock.IsDone(), true)

	// ----------- Test One Result -------------- //

	// Setup mock request.
	gock.New(apiBaseUrl).
		Get("/markets/search").
		Reply(200).
		BodyString(`{"securities":{"security":{"symbol":"SBUX","exchange":"Q","type":"stock","description":"Starbucks Corp"}}}`)

	// Make API call
	symbols2, err := tradier.SearchByCompanyName("starbucks")

	if err != nil {
		panic(err)
	}

	// Verify the data was return as expected
	st.Expect(t, len(symbols2), 1)
	st.Expect(t, symbols2[0].Name, "SBUX")
	st.Expect(t, symbols2[0].Description, "Starbucks Corp")

	// Verify that we don't have pending mocks
	st.Expect(t, gock.IsDone(), true)
}

/* End File */
