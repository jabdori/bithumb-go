# 분(Minute) 캔들

## API 정보

- **API명**: 분(Minute) 캔들
- **메서드**: GET
- **URL**: https://api.bithumb.com/v1/candles/minutes/{unit}
- **설명**: 특정 시간 간격(분)의 캔들 데이터를 조회합니다.

> 예시코드는 JavaScript, Python, JAVA에 한해서만 제공합니다.

## Request

### Path Parameters

| 필드명 | 타입 | 필수 여부 | 설명 | 기본값 |
|--------|------|----------|------|--------|
| unit | int32 | O | 분 단위. 가능한 값: 1, 3, 5, 10, 15, 30, 60, 240 | 1 |

### Query Parameters

| 필드명 | 타입 | 필수 여부 | 설명 | 기본값 |
|--------|------|----------|------|--------|
| market | string | O | 마켓 코드 (ex. KRW-BTC) | KRW-BTC |
| to | string | X | 마지막 캔들 시각 (exclusive). 비워서 요청시 가장 최근 캔들 | - |
| count | int32 | X | 캔들 개수(최대 200개까지 요청 가능) | 1 |

### Request Example

#### JavaScript

```javascript
const options = {method: 'GET', headers: {accept: 'application/json'}};

fetch('https://api.bithumb.com/v1/candles/minutes/1?market=KRW-BTC&count=1', options)
  .then(response => response.json())
  .then(response => console.log(response))
  .catch(err => console.error(err));
```

## Response

### Response Fields

| 필드명 | 타입 | 설명 |
|--------|------|------|
| market | String | 마켓명 |
| candle_date_time_utc | String | 캔들 기준 시각(UTC 기준)<br>포맷: `yyyy-MM-dd'T'HH:mm:ss` |
| candle_date_time_kst | String | 캔들 기준 시각(KST 기준)<br>포맷: `yyyy-MM-dd'T'HH:mm:ss` |
| opening_price | Double | 시가 |
| high_price | Double | 고가 |
| low_price | Double | 저가 |
| trade_price | Double | 종가 |
| timestamp | Long | 캔들 종료 시각(KST 기준) |
| candle_acc_trade_price | Double | 누적 거래 금액 |
| candle_acc_trade_volume | Double | 누적 거래량 |
| unit | Integer | 분 단위(유닛) |

### Response Codes

| 상태코드 | 설명 |
|----------|------|
| 200 | 성공 |
| 400 | 요청 파라미터 오류 |

---

*마지막 업데이트: 2026-01-01*
*API 버전: v2.1.5*
*원본 문서: https://apidocs.bithumb.com/v2.1.5/reference/분minute-캔들-1*
