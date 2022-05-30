# Donggu

![logo.png](Donggu face)

동구는 다국어 데이터 관리를 위한 가볍지만 강력한 CLI입니다.
데이터 관리부터 번역가와의 텍스트 들여오기/내보내기, 소스코드와의 간편한 연동 등 프로덕션 환경에서 요구하는 모든 상황을 대응할 수 있습니다.

## 주요 기능
### 템플릿 포맷팅

다양한 표시 형태를 지원하는 템플릿 포맷팅을 사용할 수 있습니다. 강력한 템플릿 기능을 통해 포맷팅에 필요한 코드 양을 줄일 수 있을 뿐만 아니라 번역가와 디자이너에게 텍스트의 맥락을 설명하기 쉬워집니다.

|템플릿|출력 예시|
|-|-|
|`이동거리 #{DISTANCE\|float\|6.3}km`|이동거리&nbsp;&nbsp;&nbsp;1.75km|
|`포인트 #{BALANCE\|int\|8,}원`|포인트&nbsp;&nbsp;123,456원|
|`주행 이력 #{HAS_HISTORY\|bool\|있음,없음}`|주행 이력 없음|


### 소스 코드와 쉽게 연동 가능한 라이브러리 코드 생성 
Typescript, Typescript React, Go 프로젝트에서 손쉽게 사용 가능한 라이브러리 코드를 생성할 수 있습니다. 동적으로 텍스트를 불러올 필요가 없고, 빠르면서도 타입 안정성이 보장되는 라이브러리를 명령 한줄로 생성할 수 있습니다.

```json
{
    "user.my_page.coupon_count": {
       "ko": "#{NAME}님은 쿠폰을 #{COUNT|int}개 보유중입니다.",
       "en": "Hi #{NAME}, you have #{COUNT|int} coupons."
    }
}
```
이렇게 텍스트를 저장하고 `donggu export typescript my-project`를 실행하면 아래처럼 간편한 라이브러리가 생성됩니다. 
![generated code example](docs/assets/code-generation-example.png)

### 다양한 형태로 데이터 내보내고 들어오기
동구는 JSON 형태로 번역 데이터를 관리하지만, 언제든지 CSV 혹은 HTML(준비중) 형태로 데이터를 내보내고 들여올 수 있습니다. 번역가와 디자이너에게 필요한 텍스트를 손쉽게 전달할 수 있고, 수정된 데이터도 명령 한번으로 프로젝트에 들어올 수 있습니다.
```bash
donggu export csv task.csv // 현재 프로젝트를 CSV로 내보내기
donggu merge csv done.csv  // 작업된 CSV를 현재 프로젝트에 합침
```
여기에서 어떤 데이터가 바뀌었는지의 diff 분석도 쉽게 가능합니다.

## 시작하기
[Releases](https://github.com/maasasia/donggu/releases/) 페이지에서 바이너리를 다운로드 하거나, 설치 스크립트를 이용해 다운로드 할 수 있습니다.
```bash
wget -O install.sh https://raw.githubusercontent.com/maasasia/donggu/main/install.sh && chmod +x install.sh && ./install.sh
```
바이너리가 위치한 폴더를 `PATH`에 추가해주세요.

데이터를 관리할 폴더로 이동해 새로운 프로젝트를 생성합니다.
```bash
mkdir my-project && cd my-project
donggu init
```
필요한 정보를 입력하고 나면 `metadata.json`과 `content.json`이 생성됩니다.
`metadata.json`은 메타데이터 파일로, 지원할 언어의 목록과 버전 정보등을 담습니다.
`content.json`은 번역 데이터를 담는 파일입니다.

아래에서 보다 자세한 사용 방법을 알아보세요.

## 사용방법
### 다국어 데이터의 구성과 템플릿 문법

### 코드 생성
#### Typescript

#### Typescript React

#### Go

#### 기존 프로젝트와의 연동
생성된 라이브러리는 프로젝트에 직접 추가하거나, 언어별로 지원하는 패키지 시스템을 통해 이용할 수 있습니다.
모노레포를 구성하거나 private package registry를 사용하는 등 다양한 시나리오에 대한 설명은
[기존 프로젝트와의 연동](docs/integration.md)를 참고하세요.

### 내보내기와 들여오기

### CLI 문서


## License
MIT