# 내 주문 및 체결 (MyOrder)

## Request

요청은 크게 **ticket field**, **type field**, **format field** 로 분류되며 하나의 요청에 여러 개의 **type field**를 명시할 수 있습니다. 자세한 사항은 [요청 방법 및 포맷](https://apidocs.bithumb.com/v2.1.5/reference/요청-포맷) 페이지를 확인해주시기 바랍니다.

### Type Field

수신하고 싶은 시세 정보를 나열하는 필드입니다.

| 필드명 | 타입 | 내용 | 필수 여부 | 기본 값 |
|--------|------|------|-----------|---------|
| type | String | 데이터 타입<br>- `myOrder`: 내 주문 | O | |
| codes | List | 마켓 코드 리스트<br>- 대문자로 요청해야 합니다. | X | 생략하거나 빈 배열로 요청할 경우 모든 마켓에 대한 정보를 수신합니다. |

## Response

| 필드명 | 축약형 (format :SIMPLE) | 내용 | 타입 | 값 |
|--------|------------------------|------|------|-----|
| type | ty | 타입 | String | `myOrder`: 내 주문 |
| code | cd | 마켓 코드 (ex. KRW-BTC) | String | |
| uuid | uid | 주문 고유 아이디 | String | |
| ask_bid | ab | 매수/매도 구분 | String | `ASK`: 매도, `BID`: 매수 |
| order_type | ot | 주문 타입 | String | `limit`: 지정가 주문, `price`: 시장가 주문(매수), `market`: 시장가 주문(매도) |
| state | s | 주문 상태 | String | `wait`: 체결 대기, `trade`: 체결 발생, `done`: 전체 체결 완료, `cancel`: 주문 취소 |
| trade_uuid | tuid | 체결의 고유 아이디 | String | |
| price | p | 주문 가격, 체결 가격 (state: trade 일 때) | Double | |
| volume | v | 주문량, 체결량 (state: trade 일 때) | Double | |
| remaining_volume | rv | 체결 후 남은 주문 양 | Double | |
| executed_volume | ev | 체결된 양 | Double | |
| trades_count | tc | 해당 주문에 걸린 체결 수 | Double | |
| reserved_fee | rsf | 수수료로 예약된 비용 | Double | |
| remaining_fee | rmf | 남은 수수료 | Double | |
| paid_fee | pf | 사용된 수수료 | Double | |
| executed_funds | ef | 체결된 금액 | Double | |
| trade_timestamp | ttms | 체결 타임스탬프 (millisecond) | Long | |
| order_timestamp | otms | 주문 타임스탬프 (millisecond) | Long | |
| timestamp | tms | 타임스탬프 (millisecond) | Long | |
| stream_type | st | 스트림 타입 | String | `REALTIME`: 실시간 |

## Example

### Request

**모든 마켓 정보 수신**

```json
[
  {
    "ticket": "test example"
  },
  {
    "type": "myOrder",
    "codes": []
  }
]
```

**특정 마켓 정보 수신**

```json
[
  {
    "ticket": "test example"
  },
  {
    "type": "myOrder",
    "codes": [
      "KRW-BTC"
    ]
  }
]
```

### Response

```json
{
  "type": "myOrder",
  "code": "KRW-BTC",
  "uuid": "C0101000000001818113",
  "ask_bid": "BID",
  "order_type": "limit",
  "state": "trade",
  "trade_uuid": "C0101000000001744207",
  "price": 1927000,
  "volume": 0.4697,
  "remaining_volume": 0.0803,
  "executed_volume": 0.4697,
  "trades_count": 1,
  "reserved_fee": 0,
  "remaining_fee": 0,
  "paid_fee": 0,
  "executed_funds": 905111.9,
  "trade_timestamp": 1727052318148,
  "order_timestamp": 1727052318074,
  "timestamp": 1727052318369,
  "stream_type": "REALTIME"
}
```

---

*원본 문서: [빗썸 API 레퍼런스 - 내 주문 및 체결 (MyOrder)](https://apidocs.bithumb.com/v2.1.5/reference/내-주문-및-체결-myorder)*
