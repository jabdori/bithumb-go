# 내 자산 (MyAsset)

## Request

요청은 크게 **ticket field**, **type field**, **format field** 로 분류되며 하나의 요청에 여러 개의 **type field**를 명시할 수 있습니다. 자세한 사항은 [요청 방법 및 포맷](https://apidocs.bithumb.com/v2.1.5/reference/요청-포맷) 페이지를 확인해주시기 바랍니다.

### Type Field

수신하고 싶은 시세 정보를 나열하는 필드입니다.

| 필드명 | 타입 | 내용 | 필수 여부 | 기본 값 |
|--------|------|------|-----------|---------|
| type | String | 데이터 타입<br>- `myAsset`: 내 자산 | O | |

## Response

| 필드명 | 축약형 (format :SIMPLE) | 내용 | 타입 | 값 |
|--------|------------------------|------|------|-----|
| type | ty | 타입 | String | `myAsset`: 내 자산 |
| assets | ast | 자산 리스트 | List of Objects | |
| assets.currency | ast.cu | 화폐를 의미하는 영문 대문자 코드 | String | |
| assets.balance | ast.b | 주문가능 수량 | Double | |
| assets.locked | ast.l | 주문 중 묶여있는 수량 | Double | |
| asset_timestamp | asttms | 자산 타임스탬프 (millisecond) | Long | |
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
    "type": "myAsset"
  }
]
```

### Response

```json
{
  "type": "myAsset",
  "assets": [
    {
      "currency": "KRW",
      "balance": "2061832.35",
      "locked": "3824127.3"
    }
  ],
  "asset_timestamp": 1727052537592,
  "timestamp": 1727052537687,
  "stream_type": "REALTIME"
}
{
  "type": "myAsset",
  "assets": [
    {
      "currency": "BTC",
      "balance": "156.70564833",
      "locked": "38.81945789"
    }
  ],
  "asset_timestamp": 1727052537592,
  "timestamp": 1727052537690,
  "stream_type": "REALTIME"
}
```

---

*원본 문서: [빗썸 API 레퍼런스 - 내 자산 (MyAsset)](https://apidocs.bithumb.com/v2.1.5/reference/내-자산-myasset)*
