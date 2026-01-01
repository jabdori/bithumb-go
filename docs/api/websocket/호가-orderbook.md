# 호가 (Orderbook)

## Request

요청은 크게 **ticket field**, **type field**, **format field** 로 분류되며 하나의 요청에 여러 개의 **type field**를 명시할 수 있습니다. 자세한 사항은 [요청 방법 및 포맷](https://apidocs.bithumb.com/v2.1.5/reference/요청-포맷) 페이지를 확인해주시기 바랍니다.

### Type Field

수신하고 싶은 시세 정보를 나열하는 필드입니다. `is_only_snapshot`, `is_only_realtime` 필드는 생략 가능하며 모두 생략할 경우 스냅샷과 실시간 데이터 둘 다 수신합니다.

| 필드명 | 타입 | 내용 | 필수 여부 | 기본 값 |
|--------|------|------|-----------|---------|
| type | String | 데이터 타입<br>- `orderbook`: 호가 | O | |
| codes | List | 마켓 코드 리스트<br>- 대문자로 요청해야 합니다. | O | |
| level | Double | 모아보기 단위 | X | 1 |
| isOnlySnapshot | Boolean | 스냅샷 시세만 제공 | X | false |
| isOnlyRealtime | Boolean | 실시간 시세만 제공 | X | false |

## Response

| 필드명 | 축약형 (format :SIMPLE) | 내용 | 타입 | 값 |
|--------|------------------------|------|------|-----|
| type | ty | 타입 | String | `orderbook`: 호가 |
| code | cd | 마켓 코드 (ex. KRW-BTC) | String | |
| total_ask_size | tas | 호가 매도 총 잔량 | Double | |
| total_bid_size | tbs | 호가 매수 총 잔량 | Double | |
| orderbook_units | obu | 호가 | List of Objects | |
| orderbook_units.ask_price | obu.ap | 매도 호가 | Double | |
| orderbook_units.bid_price | obu.bp | 매수 호가 | Double | |
| orderbook_units.ask_size | obu.as | 매도 잔량 | Double | |
| orderbook_units.bid_size | obu.bs | 매수 잔량 | Double | |
| timestamp | tms | 타임스탬프 (millisecond) | Long | |
| level | lv | 호가 모아보기 단위 (default: 1, 기본 호가단위) | Double | 모아보기 단위 |

## Example

### Request

- `level` 값은 필수가 아니며, 제외될 경우 DEFAULT(1), 기본 호가단위로 내려갑니다.

```json
[
  {
    "ticket": "test example"
  },
  {
    "type": "orderbook",
    "codes": [
      "KRW-BTC",
      "KRW-ETH.3"
    ],
    "level": 10
  },
  {
    "format": "DEFAULT"
  }
]
```

종목별로 각기 다른 모아보기 `level` 값을 지정하기 위해서는 아래와 같이 요청할 수 있습니다.

```json
[
  {
    "ticket": "test example"
  },
  {
    "type": "orderbook",
    "codes": [
      "KRW-BTC"
    ],
    "level": 1000
  },
  {
    "type": "orderbook",
    "codes": [
      "KRW-XRP"
    ],
    "level": 1
  },
  {
    "format": "DEFAULT"
  }
]
```

### Response

```json
{
  "type": "orderbook",
  "code": "KRW-BTC",
  "total_ask_size": 450.3526,
  "total_bid_size": 63.3006,
  "orderbook_units": [
    {
      "ask_price": 478800,
      "bid_price": 478300,
      "ask_size": 4.3478,
      "bid_size": 5.6370
    },
    {
      "ask_price": 489700,
      "bid_price": 477900,
      "ask_size": 2.3642,
      "bid_size": 0.9705
    },
    {
      "ask_price": 493100,
      "bid_price": 471200,
      "ask_size": 411.8686,
      "bid_size": 3.9279
    },
    {
      "ask_price": 493300,
      "bid_price": 471100,
      "ask_size": 2.0241,
      "bid_size": 1.4699
    },
    {
      "ask_price": 493700,
      "bid_price": 471000,
      "ask_size": 1.7870,
      "bid_size": 2.2573
    },
    {
      "ask_price": 493800,
      "bid_price": 470700,
      "ask_size": 3.9372,
      "bid_size": 9.7805
    },
    {
      "ask_price": 494900,
      "bid_price": 470400,
      "ask_size": 5.7560,
      "bid_size": 0.8093
    },
    {
      "ask_price": 495300,
      "bid_price": 470300,
      "ask_size": 3.6418,
      "bid_size": 4.6606
    },
    {
      "ask_price": 495700,
      "bid_price": 470100,
      "ask_size": 2.9617,
      "bid_size": 5.4907
    },
    {
      "ask_price": 495800,
      "bid_price": 469700,
      "ask_size": 0.2349,
      "bid_size": 2.3941
    },
    {
      "ask_price": 496100,
      "bid_price": 469600,
      "ask_size": 2.6019,
      "bid_size": 4.5505
    },
    {
      "ask_price": 496800,
      "bid_price": 469500,
      "ask_size": 3.4651,
      "bid_size": 5.4469
    },
    {
      "ask_price": 496900,
      "bid_price": 469200,
      "ask_size": 0.8400,
      "bid_size": 10.1685
    },
    {
      "ask_price": 497400,
      "bid_price": 469100,
      "ask_size": 2.1924,
      "bid_size": 5.1646
    },
    {
      "ask_price": 497900,
      "bid_price": 469000,
      "ask_size": 2.3299,
      "bid_size": 0.5723
    }
  ],
  "level": 1,
  "timestamp": 1725930007672,
  "stream_type": "REALTIME"
}
```

---

*원본 문서: [빗썸 API 레퍼런스 - 호가 (Orderbook)](https://apidocs.bithumb.com/v2.1.5/reference/호가-orderbook)*
