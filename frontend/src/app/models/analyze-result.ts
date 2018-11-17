//
// Date: 11/16/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

//
// AnalyzeResult
//
export class AnalyzeResult 
{
  UnderlyingPrice: number;
  Profit: number;

  //
  // Build from JSON list.
  //
  fromJsonList(json: Object[]): AnalyzeResult[] 
  {
    let result = [];

    if (!json) 
    {
      return result;
    }

    for (let i = 0; i < json.length; i++) 
    {
      result.push(new AnalyzeResult().fromJson(json[i]));
    }

    // Return happy
    return result;
  }

  //
  // Json to Object.
  //
  fromJson(json: Object): AnalyzeResult 
  {
    let obj = new AnalyzeResult();

    obj.Profit = json["profit"];
    obj.UnderlyingPrice = json["underlying_price"];

    return obj;
  }  
}

export class AnalyzeLeg
{
  SymbolStr: string;
  Qty: number;
}

/* End File */