# 일(Day) 캔들

## API 정보

- **API명**: 일(Day) 캔들
- **메서드**: GET
- **URL**: https://api.bithumb.com/v1/candles/days
- **설명**: 일별 캔들 데이터를 조회합니다.

> 예시코드는 JavaScript, Python, JAVA에 한해서만 제공합니다.

## Request

### Query Parameters

| 필드명 | 타입 | 필수 여부 | 설명 | 기본값 |
|--------|------|----------|------|--------|
| market | string | O | 마켓 코드 (ex. KRW-BTC) | - |
| to | string | X | 마지막 캔들 시각 (exclusive). 비워서 요청시 가장 최근 캔들 | - |
| count | int32 | X | 캔들 개수(최대 200개까지 요청 가능) | 1 |
| convertingPriceUnit | string | X | 종가 환산 화폐 단위 (생략 가능, KRW로 명시할 시 원화 환산 가격을 반환) | - |

### Request Example

#### JavaScript

```javascript
const options = {method: 'GET', headers: {accept: 'application/json'}};

fetch('https://api.bithumb.com/v1/candles/days?count=1', options)
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
| prev_closing_price | Double | 전일 종가(UTC 0시 기준) |
| change_price | Double | 전일 종가 대비 변화 금액 |
| change_rate | Double | 전일 종가 대비 변화량 |
| converted_trade_price | Double | 종가 환산 화폐 단위로 환산된 가격 (요청에 `convertingPriceUnit` 파라미터가 없는 경우 해당 필드는 반환되지 않음) |

> **참고**: `convertingPriceUnit` 파라미터의 경우, 원화 마켓이 아닌 다른 마켓(ex. BTC)의 일봉 요청시 종가를 명시된 파라미터 값으로 환산해 `converted_trade_price` 필드에 추가하여 반환합니다. 현재는 원화(`KRW`)로 변환하는 기능만 제공하며 추후 기능을 확장할 수 있습니다.

### Response Codes

| 상태코드 | 설명 |
|----------|------|
| 200 | 성공 |
| 400 | 요청 파라미터 오류 |

---

*마지막 업데이트: 2026-01-01*
*API 버전: v2.1.5*
*원본 문서: https://apidocs.bithumb.com/v2.1.5/reference/일day-캔들*
