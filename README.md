# 프로젝트 개요

ZeroMQ의 Go 언어 라이브러리 선택지는 두 가지가 있습니다.

1. **goczmq**
   
   goczmq는 ZeroMQ의 Go 언어 바인딩 라이브러리입니다. ZeroMQ의 CZMQ C 라이브러리를 Go에서 사용할 수 있도록 만들어졌습니다.
   
   - 고수준 API 제공: 사용자가 쉽게 이해하고 사용할 수 있는 API를 제공합니다.
   - CZMQ 기반: CZMQ 라이브러리의 기능과 안정성을 활용할 수 있습니다.
   - 편리한 메모리 관리: C 언어의 메모리 관리를 Go 스타일로 추상화하여 제공합니다.

   **=> 간단한 사용: goczmq (고수준 API, 편리한 메모리 관리)**

2. **zmq4**
   
   zmq4는 ZeroMQ의 Go 언어 바인딩 라이브러리로, ZeroMQ의 원시 C API를 Go 언어로 직접 바인딩한 것입니다.
   
   - 저수준 API 제공: ZeroMQ의 원시 API에 직접 접근하여 더 세밀한 제어가 가능합니다.
   - 성능 최적화: 저수준 API를 사용하여 성능 최적화가 가능합니다.
   - ZeroMQ의 최신 기능 지원: ZeroMQ의 최신 기능을 빠르게 반영할 수 있습니다.

   **=> 성능과 제어: zmq4 (저수준 API, 성능 최적화)**

## 선택 기준

- **단순함과 사용 편의성** -> goczmq
- **성능과 세밀한 제어** -> zmq4
  
   **다양한 통신 패턴과 비동기 메시징 기능을 제공하는 libzmq를 바인딩한 zmq4를 사용하기로 결정**   
   **GO 언어의 성능과 ZeroMQ의 다양한 통신 패턴을 원시 API로 직접 접근할 수 있는 zmq4를 사용하기로 결정**   
   **=> 통신 속도의 최적화 및 성능 향상 기대**   


## 프로젝트 진행:
1. Go 언어 설치 및 프로젝트 환경 설정
   - Homebrew를 통해 Go 설치
   - go tool 패키지 설치
   - VSCode에 Go Extension 설치
2. ZeroMQ를 사용하기 위해 Go에서 `github.com/pebbe/zmq4` 패키지를 설치합니다.
   ```sh
   go get -u github.com/pebbe/zmq4
