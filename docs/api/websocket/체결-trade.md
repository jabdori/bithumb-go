# 체결 (Trade)

## Request

요청은 크게 **ticket field**, **type field**, **format field** 로 분류되며 하나의 요청에 여러 개의 **type field**를 명시할 수 있습니다. 자세한 사항은 [요청 방법 및 포맷](https://apidocs.bithumb.com/v2.1.5/reference/요청-포맷) 페이지를 확인해주시기 바랍니다.

### Type Field

수신하고 싶은 시세 정보를 나열하는 필드입니다. `is_only_snapshot`, `is_only_realtime` 필드는 생략 가능하며 모두 생략할 경우 스냅샷과 실시간 데이터 둘 다 수신합니다.

| 필드명 | 타입 | 내용 | 필수 여부 | 기본 값 |
|--------|------|------|-----------|---------|
| type | String | 데이터 타입<br>- `trade`: 체결 | O | |
| codes | List | 마켓 코드 리스트<br>- 대문자로 요청해야 합니다. | O | |
| isOnlySnapshot | Boolean | 스냅샷 시세만 제공 | X | false |
| isOnlyRealtime | Boolean | 실시간 시세만 제공 | X | false |

## Response

| 필드명 | 축약형 (format :SIMPLE) | 내용 | 타입 | 값 |
|--------|------------------------|------|------|-----|
| type | ty | 타입 | String | `trade`: 체결 |
| code | cd | 마켓 코드 (ex. KRW-BTC) | String | |
| trade_price | tp | 체결 가격 | Double | |
| trade_volume | tv | 체결량 | Double | |
| ask_bid | ab | 매수/매도 구분 | String | `ASK`: 매도, `BID`: 매수 |
| prev_closing_price | pcp | 전일 종가 | Double | |
| change | c | 전일 대비 | String | `RISE`: 상승, `EVEN`: 보합, `FALL`: 하락 |
| change_price | cp | 부호 없는 전일 대비 값 | Double | |
| trade_date | tdt | 최근 거래 일자(KST) | String | yyyy-MM-dd |
| trade_time | ttm | 최근 거래 시각(KST) | String | HH:mm:ss |
| trade_timestamp | ttms | 체결 타임스탬프 (milliseconds) | Long | |
| timestamp | tms | 타임스탬프 (millisecond) | Long | |
| sequential_id | sid | 체결 번호 (Unique) | Long | |
| stream_type | st | 스트림 타입 | String | `SNAPSHOT`: 스냅샷, `REALTIME`: 실시간 |

> **참고**: `sequential_id` 필드는 체결의 유일성을 판단하기 위한 근거로 쓰일 수 있습니다. 하지만 체결 순서를 보장하지는 못합니다.

## Example

### Request

- `KRW-BTC`, `KRW-ETH`

```json
[
  {
    "ticket": "test example"
  },
  {
    "type": "trade",
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
  "type": "trade",
  "code": "KRW-BTC",
  "trade_price": 489700,
  "trade_volume": 1.4825,
  "ask_bid": "BID",
  "prev_closing_price": 484500,
  "change": "RISE",
  "change_price": 5200,
  "trade_date": "2024-09-10",
  "trade_time": "09:58:54",
  "trade_timestamp": 1725929934373,
  "sequential_id": 17259299343730000,
  "timestamp": 1725929934483,
  "stream_type": "REALTIME"
}
```

---

*원본 문서: [빗썸 API 레퍼런스 - 체결 (Trade)](https://apidocs.bithumb.com/v2.1.5/reference/체결-trade)*
