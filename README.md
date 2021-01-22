#### This project is aimed to simplify creating Go Web Applications by bundling logging, JWT authentication, MongoDB, Redis, and Kafka integration.
</br>

---
</br>

## Project Overview

This project contains 6 main packages:

1. Server
2. App
3. API
4. Conf
5. Mock
6. Model

### Server

Server separates the http server, external system integration and auth logic from business logic implementation

### App
App package should contain all the code to implement business logic.

### API
API exposes the business logic through API endpoints. API package should implement CRUD operations & validations of the endpoints.

### Conf
Conf only contains 2 main files `default.toml` used by server to fetch the project configuration and `test.toml` used for unit testing.

### Mock
Mock should contain all the mocked interfaces and dependencies. GoMock is used as a tool for generating and mocking the interfaces.

### Model
Model consists all database related schemas for struct|JSON|BSON encoding and decoding.

---
</br>

## How to start with the project
</br>

### Create an endpoint

* go to `api>routes.go` file

    **to create an endpoint**

        a.Router.Root.Handle(
            "/", 
            a.requestHandler(<controller func>)
        ).Methods("GET")

    
    **to create an endpoint that start with /api/**

        a.Router.APIRoot.Handle(
            "/", 
            a.requestHandler(<controller func>)
        ).Methods("POST")


    **to create an endpoint for static files**

        a.Router.StaticRoot.Handle(
            "/", 
            a.requestHandler(<controller func>)
        ).Methods("GET")

</br>

### Create an endpoint controller func

* For modularisation purpose, create a separate go file to similar handler functions

* Function should be implemented like : 

        func (a *API) funcName(
            requestCTX *handler.RequestContext, 
            w http.ResponseWriter, 
            r *http.Request
        ) { 
            Your Code 
        }  

</br>

### Create a new service

* Create a new go file inside `app` package (Why? -> code separation and better visibility)

* Create an interface with service name `type Example interface` and define all the methods that the service should implement.

        type Example interface {
            SayHello(bool) (string, error)
            SaveHello(*SaveHelloOpts) (*SaveHelloResp, error)
        }

* Create a `ServiceOpts` that takes all the external dependencies and `app` instance

        type ExampleOpts struct {
            App    *App
            DB     *mongo.Database
            Logger *zerolog.Logger
        }


* Implement all service interface methods

        func (e *ExampleImpl) SaveHello(opts *SaveHelloOpts) (*SaveHelloResp, error) {
            res, err := e.DB.Collection("hello").InsertOne(context.TODO(), opts)
            if err != nil {
                return nil, err
            }
            out := SaveHelloResp{
                Name: opts.Name,
                ID:   res.InsertedID.(primitive.ObjectID),
            }
            return &out, nil
        }

        func (e *ExampleImpl) SayHello(wantErr bool) (string, error) {
            if wantErr {
                return "", errors.New("hello error")
            }
            return "Hello", nil
        }


* Define a constructor function with name `Init<Service>` to return the Service implementation instance of the service


        func InitExample(opts *ExampleOpts) Example {
            e := &ExampleImpl{
                App:    opts.App,
                DB:     opts.DB,
                Logger: opts.Logger,
            }
            return e
        }

</br>

---
</br>

## How to write unit test

Go has an idiomatic way writing test cases which is table driven test cases. Thus, this doc covers how to write table driven unit test cases.

</br>

### Unit test for endpoint and controller func
</br>

An example of how table driven struct should be defined


    // Setting up API
    api := NewTestAPI(getTestConfig())
    .
    .
    .
    tests := []struct {
		name          string
		url           string
		method        string
		body          io.Reader
		saveHelloOpts *app.SaveHelloOpts
		buildStubs    func(ex *mock.MockExample)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	} { { ... }, { ... } }

    for _, tt := range tests {
            t.Run(tt.name, func(t *testing.T) {
                ctrl := gomock.NewController(t)
                defer ctrl.Finish()

                // building stub
                ex := mock.NewMockExample(ctrl)
                tt.buildStubs(ex)
                // setting stub
                api.App.Example = ex

                recorder := httptest.NewRecorder()
                req, err := http.NewRequest(tt.method, tt.url, tt.body)
                assert.Nil(t, err)
                api.Router.Root.ServeHTTP(recorder, req)
                tt.checkResponse(t, recorder)
            })
        }


### Unit test for app service

    // Setting Up APP
    app := NewTestApp(getTestConfig())
    defer CleanTestApp(app)
    type fields struct {
        App    *App
        DB     *mongo.Database
        Logger *zerolog.Logger
    }
    type args struct {
        wantErr bool
    }
    tests := []struct {
        name    string
        args    args
        fields  fields
        want    string
        wantErr bool
    } { { ... }, { ... } }

    for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ExampleImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			got, err := e.SayHello(tt.args.wantErr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExampleImpl.SayHello() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExampleImpl.SayHello() = %v, want %v", got, tt.want)
			}
		})
	}


---
## Generate service mock

`//go:generate $GOPATH/bin/mockgen -destination=../mock/mock_example.go -package=mock go-app/app Example`

Modify the line above accordingly and at the top of the file where service interface is defined.

---
## Run and Build

***to run the project*** -->    `go run main.go`

***to build project*** --->     `go build main.go`