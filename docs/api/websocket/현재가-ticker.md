# 현재가 (Ticker)

## Request

요청은 크게 **ticket field**, **type field**, **format field** 로 분류되며 하나의 요청에 여러 개의 **type field**를 명시할 수 있습니다. 자세한 사항은 [요청 방법 및 포맷](https://apidocs.bithumb.com/v2.1.5/reference/요청-포맷) 페이지를 확인해주시기 바랍니다.

### Type Field

수신하고 싶은 시세 정보를 나열하는 필드입니다. `is_only_snapshot`, `is_only_realtime` 필드는 생략할 수 있으며 둘 다 생략할 경우 스냅샷과 실시간 데이터 모두 수신합니다.

| 필드명 | 타입 | 내용 | 필수 여부 | 기본 값 |
|--------|------|------|-----------|---------|
| type | String | 데이터 타입<br>- `ticker`: 현재가 | O | |
| codes | List | 마켓 코드 리스트<br>- 대문자로 요청해야 합니다. | O | |
| isOnlySnapshot | Boolean | 스냅샷 시세만 제공 | X | false |
| isOnlyRealtime | Boolean | 실시간 시세만 제공 | X | false |

## Response

| 필드명 | 축약형 (format :SIMPLE) | 내용 | 타입 | 값 |
|--------|------------------------|------|------|-----|
| type | ty | 타입 | String | `ticker`: 현재가 |
| code | cd | 마켓 코드 (ex. KRW-BTC) | String | |
| opening_price | op | 시가 | Double | |
| high_price | hp | 고가 | Double | |
| low_price | lp | 저가 | Double | |
| trade_price | tp | 현재가 | Double | |
| prev_closing_price | pcp | 전일 종가 | Double | |
| change | c | 전일 대비 | String | `RISE`: 상승, `EVEN`: 보합, `FALL`: 하락 |
| change_price | cp | 부호 없는 전일 대비 값 | Double | |
| signed_change_price | scp | 전일 대비 값 | Double | |
| change_rate | cr | 부호 없는 전일 대비 등락율 | Double | |
| signed_change_rate | scr | 전일 대비 등락율 | Double | |
| trade_volume | tv | 가장 최근 거래량 | Double | |
| acc_trade_volume | atv | 누적 거래량(KST 0시 기준) | Double | |
| acc_trade_volume_24h | atv24h | 24시간 누적 거래량 | Double | |
| acc_trade_price | atp | 누적 거래대금(KST 0시 기준) | Double | |
| acc_trade_price_24h | atp24h | 24시간 누적 거래대금 | Double | |
| trade_date | tdt | 최근 거래 일자(KST) | String | yyyyMMdd |
| trade_time | ttm | 최근 거래 시각(KST) | String | HHmmss |
| trade_timestamp | ttms | 체결 타임스탬프 (milliseconds) | Long | |
| ask_bid | ab | 매수/매도 구분 | String | `ASK`: 매도, `BID`: 매수 |
| acc_ask_volume | aav | 누적 매도량 | Double | |
| acc_bid_volume | abv | 누적 매수량 | Double | |
| highest_52_week_price | h52wp | 52주 최고가 | Double | |
| highest_52_week_date | h52wdt | 52주 최고가 달성일 | String | yyyy-MM-dd |
| lowest_52_week_price | l52wp | 52주 최저가 | Double | |
| lowest_52_week_date | l52wdt | 52주 최저가 달성일 | String | yyyy-MM-dd |
| market_state | ms | 거래상태 | String | |
| is_trading_suspended | its | 거래 정지 여부 | Boolean | |
| delisting_date | dd | 거래지원 종료일 | Date | |
| market_warning | mw | 유의 종목 여부 | String | `NONE`: 해당없음, `CAUTION`: 거래유의 |
| timestamp | tms | 타임스탬프 (millisecond) | Long | |
| stream_type | st | 스트림 타입 | String | `SNAPSHOT`: 스냅샷, `REALTIME`: 실시간 |

## Example

### Request

- `KRW-BTC`, `KRW-ETH`

```json
[
  {
    "ticket": "test example"
  },
  {
    "type": "ticker",
    "codes": [
      "KRW-BTC",
      "KRW-ETH"
    ]
  },
  {
    "format": "DEFAULT"
  }
]
```

### Response

```json
{
  "type": "ticker",
  "code": "KRW-BTC",
  "opening_price": 484500,
  "high_price": 493100,
  "low_price": 472500,
  "trade_price": 493100,
  "prev_closing_price": 484500,
  "change": "RISE",
  "change_price": 8600,
  "signed_change_price": 8600,
  "change_rate": 0.01775026,
  "signed_change_rate": 0.01775026,
  "trade_volume": 1.2567,
  "acc_trade_volume": 225.622,
  "acc_trade_volume_24h": 13386.15417512,
  "acc_trade_price": 108663718.238256,
  "acc_trade_price_24h": 8230696760.346009,
  "trade_date": "20240910",
  "trade_time": "091617",
  "trade_timestamp": 1725927377820,
  "ask_bid": "BID",
  "acc_ask_volume": 106.7561,
  "acc_bid_volume": 118.8659,
  "highest_52_week_price": 999999000,
  "highest_52_week_date": "2024-06-18",
  "lowest_52_week_price": 1000,
  "lowest_52_week_date": "2024-06-18",
  "market_state": "ACTIVE",
  "is_trading_suspended": false,
  "delisting_date": null,
  "market_warning": "NONE",
  "timestamp": 1725927377931,
  "stream_type": "REALTIME"
}
```

---

*원본 문서: [빗썸 API 레퍼런스 - 현재가 (Ticker)](https://apidocs.bithumb.com/v2.1.5/reference/현재가-ticker)*
