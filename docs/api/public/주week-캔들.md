# 주(Week) 캔들

## API 정보

- **API명**: 주(Week) 캔들
- **메서드**: GET
- **URL**: https://api.bithumb.com/v1/candles/weeks
- **설명**: 주별 캔들 데이터를 조회합니다.

> 예시코드는 JavaScript, Python, JAVA에 한해서만 제공합니다.

## Request

### Query Parameters

| 필드명 | 타입 | 필수 여부 | 설명 | 기본값 |
|--------|------|----------|------|--------|
| market | string | O | 마켓 코드 (ex. KRW-BTC) | - |
| to | string | X | 마지막 캔들 시각 (exclusive). 비워서 요청시 가장 최근 캔들 | - |
| count | int32 | X | 캔들 개수(최대 200개까지 요청 가능) | 1 |

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
| first_day_of_period | String | 캔들 기간의 가장 첫 날 |

### Response Codes

| 상태코드 | 설명 |
|----------|------|
| 200 | 성공 |
| 400 | 요청 파라미터 오류 |

---

*마지막 업데이트: 2026-01-01*
*API 버전: v2.1.5*
*원본 문서: https://apidocs.bithumb.com/v2.1.5/reference/주week-캔들*
