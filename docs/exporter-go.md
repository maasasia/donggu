# Typescript React Exporter
다국어 데이터를 사용할 수 있는 Typescript React 패키지를 만들어 주는 exporter입니다.

기본적인 기능은 [Typescript Exporter](exporter-ts.md) 동일하지만,
리액트에서 편리한 기능이 몇가지 추가됩니다.

- [사용 방법](#usage)
- [연동 가이드](#integration)

## 메타데이터 파일 준비
`metadata.json`에서 `exporter_options`의 `ts-react` 키 아래 옵션을 지정할 수 있습니다.
다음과 같은 정보가 필요합니다.
- `package_name`: 생성되는 라이브러리의 npm 패키지명

## 사용 방법 <span id="usage"></span>


### Donggu Struct
다국어 데이터는 Donggu 구조체를 통해 접근할 수 있습니다.

Donggu 구조체는 allocation 횟수를 최소화 할 수 있는 구조로 (출력 1회당 0 혹은 1회) 설계되어 있어
성능에 대한 우려 없이 사용할 수 있습니다.

Donggu 생성자는 **텍스트에 사용할 언어를 판단하는** 함수를 필요로 합니다.

```go

```

### 함수 사용
다국어 데이터들은 키만 `camelCase`로 변경되어 그대로 사용할 수 있습니다.
```json
{
    "screens.login_page.modal.success": {
        "ko": "로그인 성공"
    },
    "screens.login_page.ban_message": {
        "ko": "#{USER_NAME|string}님은 #{MONTH|int|02}월 #{DAY|int|02}일까지 제한된 사용자입니다."
    }
}
```
이런 데이터 항목들이 있다면 
```go
text := donggu.screens().loginPage().modal().success()
fmt.Println(text)
// 로그인 성공
```
함수를 호출하여 텍스트를 가져올 수 있습니다. 언어는 Donggu 인스턴스를 만들때 넘겨주었던 선호 언어 판단 함수의 결과로 결정됩니다.

템플릿들은 함수의 매개변수로 넘겨주며 템플릿 자료형에 지정한 타입을 그대로 사용합니다.
매개변수 순서는 각 함수의 시그니쳐를 참고해주세요.

```go
text := donggu.screens().loginPage().banMessage("홍길동", 8, 31)
fmt.Println(text)
// 홍길동님은 08월 31일까지 제한된 사용자입니다.
```

## 연동 가이드 <span id="integration"></span>
`작성중`

### [Gin](https://github.com/gin-gonic/gin)과 연동하기

### [text/language](https://pkg.go.dev/golang.org/x/text/language#hdr-Matching_preferred_against_supported_languages) Matcher 사용하기