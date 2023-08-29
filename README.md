[English](README.md) | [中文版](README_cn.md)

[![License: Apache 2](https://img.shields.io/badge/license-Apache%202-green)](https://github.com/callmegema/tuya-connector-go/blob/master/LICENSE "License")
![Version: 1.0.0](https://img.shields.io/badge/version-1.0.0-blue)

The `tuya-connector-go` framework maps cloud APIs to local APIs based on simple configurations and flexible extension mechanisms. You can subscribe to the distribution of cloud messages as local events. You can put all the focus on business logic without taking care of server-side programming nor relational databases. The OpenAPI or message subscription process is simplified, so you can focus on service logic and promote development efficiency.

## Quick start

#### Preparation

-   `AccessId` & `AccessKey`: the authorization key.
-   API address: Select the API address that is supported by the areas in which your business is deployed.
-   Pulsar address: Select the Pulsar address that is supported by business areas in which your business is deployed.

#### Configuration

-   Method 1: recommended. Set environment variables, and read the configuration from the environment variables when the project starts.<br/>
    export TUYA_API_HOST=https://xxxxx.com<br/>
    export TUYA_ACCESSID=xxxxxx<br/>
    export TUYA_ACCESSKEY=xxxxxxx<br/>
    export TUYA_MSG_HOST=pulsar+ssl://xxxxxx

-   Method 2: Define variables in the project, with hard-coded settings.

#### Usage

##### 1. Initialize the application (default configuration)

```go
func main() {
   // Default configuration
   connector.InitWithOptions()
   // Start the message subscription service
   go messaging.Listener()
   // Start the API service
   r := router.NewGinEngin()
   go r.Run("0.0.0.0:2021")
   watitSignal()
}

func watitSignal() {
   quitCh := make(chan os.Signal, 1)
   signal.Notify(quitCh, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
   for {
      select {
      case c := <-quitCh:
         extension.GetMessage(constant.TUYA_MESSAGE).Stop()
         logger.Log.Infof("receive sig:%v, shutdown the http server...", c.String())
         return
      }
   }
}
```

Alternatively, use a custom configuration to initialize the application:

```go
func main() {
   // Custom configuration
   connector.InitWithOptions(env.WithApiHost("xxxx"),
      env.WithMsgHost("xxxx"),
      env.WithAccessID("xxxx"),
      env.WithAccessKey("xxxx"))
   // Start the service
   go messaging.Listener()
   r := router.NewGinEngin()
   go r.Run("0.0.0.0:2021")
   watitSignal()
}
```

##### 2. Call OpenAPI (taking Gin framework for example)

```go
// Use Gin to create a route
func initRouter(r *gin.Engine) {
   r.GET("/devices/:device_id", GetDevice)
}

func GetDevice(c *gin.Context) {
   device_id := c.Param("device_id")
   resp := &GetDeviceResponse{}
   // Initiate an API request
   err := connector.MakeGetRequest(
      context.Background(),
      connector.WithAPIUri(fmt.Sprintf("/v1.0/devices/%s", device_id)),
      connector.WithResp(resp))
   if err != nil {
      fmt.Println("err:", err.Error())
      c.Abort()
      return
   }
   c.JSON(200, resp)
}

// Data structure returned by OpenAPI
type GetDeviceResponse struct {
   Code    int         `json:"code"`
   Msg     string      `json:"msg"`
   Success bool        `json:"success"`
   Result  interface{} `json:"result"`
   T       int64       `json:"t"`
}
```

When you call OpenAPI, you can customize the processing logic according to different error codes. You need to create the `struct` class that implements the `IError` interface.

```go
type DeviceError struct {
}

func (d *DeviceError) Process(ctx context.Context, code int, msg string) {
   logger.Log.Error(code, msg)
}
```

```go
func GetDevice(c *gin.Context) {
   device_id := c.Param("device_id")
   resp := new(map[string]interface{})
   err := connector.MakeGetRequest(
      context.Background(),
      connector.WithAPIUri(fmt.Sprintf("/v1.0/devices/%s", device_id)),
      connector.WithResp(resp),
      // Set custom event handling according to error codes
      connector.WithErrProc(1102, &DeviceError{}))
   if err != nil {
      fmt.Println("err:", err.Error())
      c.Abort()
      return
   }
   c.JSON(200, resp)
}
```

##### 3. Subscribe to message events

```go
func Listener() {
   // Initialize the message queue client
   extension.GetMessage(constant.TUYA_MESSAGE).InitMessageClient()
   // The message event that the specified device changes the name
   extension.GetMessage(constant.TUYA_MESSAGE).SubEventMessage(func(m *event.NameUpdateMessage) {
      logger.Log.Info("=========== name update： ==========")
      logger.Log.Info(m)
   })
   // The message event that the specified device reports data
   extension.GetMessage(constant.TUYA_MESSAGE).SubEventMessage(func(m *event.StatusReportMessage) {
      logger.Log.Info("=========== report data： ==========")
      for _, v := range m.Status {
         logger.Log.Infof(v.Code, v.Value)
      }
   })
}
```

## Custom extension implementation

-   `IError` interface

    When an error code is returned after you request OpenAPI, you can customize the struct that implements `IError` to deal with different error codes.<br/>

-   `IToken` interface

    Support the token struct to customize `IToken`, manage the token lifecycle, get or refresh the token, and locally cache the token information. The object is injected into the underlying framework when the service is initialized.<br/>

    > ```
    > extension.GetToken(constant.CUSTOM_TOKEN)
    > ```

-   `IHeader` interface

    Support the custom header struct. Customize processing logic when you request Tuya IoT Cloud OpenAPI, including attribute values and signatures that need to be added. The object is injected into the underlying framework when the service is initialized.<br/>

    > ```
    > extension.GetHeader(constant.CUSTOM_HEADER)
    > ```

-   `ISign` interface

    Customize the signature logic with the ISign interface. The object is injected into the underlying framework when the service is initialized.<br/>

    > ```
    > extension.GetSign(constant.CUSTOM_SIGN)
    > ```

-   `IEventMessage` interface

    Customize the message event subscription. The connection methods, receiving messages, and data decryption may be different for message queues. In this case, you need to customize the message subscription logic. The object is injected into the underlying framework when the service is initialized.<br/>

    > ```
    > extension.GetMessage(constant.CUSTOM_MESSAGE)
    > ```

### Message subscription

You can add the processing function for required events. The framework includes all Tuya's cloud message event types. The message data contains ciphertext messages and plaintext messages.

| **Message event**               | **BizCode**              | **Description**                    |
| ------------------------------- | ------------------------ | ---------------------------------- |
| StatusReportMessage             | statusReport             | Report data to the cloud.          |
| OnlineMessage                   | online                   | A device is online.                |
| OfflineMessage                  | offline                  | A device is offline.               |
| NameUpdateMessage               | nameUpdate               | Modify the device name.            |
| DpNameUpdateMessage             | dpNameUpdate             | Modify the name of a data point.   |
| DeleteMessage                   | delete                   | Remove a device.                   |
| BindUserMessage                 | bindUser                 | Bind the device to a user account. |
| UpgradeStatusMessage            | upgradeStatus            | The update status.                 |
| AutomationExternalActionMessage | automationExternalAction | Automate an external action.       |
| SceneExecuteMessage             | sceneExecute             | Execute a scene.                   |

### Architecture of the framework

![Integration of the framework](https://github.com/callmegema/tuya-connector-go/blob/fetch-dev/assets/architect.jpg)
![Integration and extensions](https://github.com/callmegema/tuya-connector-go/blob/fetch-dev/assets/integration%26extension.jpg)

### Core module design

![Design description](https://github.com/callmegema/tuya-connector-go/blob/fetch-dev/assets/tuya-connector-go.png)

## Technical support

You can get technical support from Tuya in the following services:

-   Help Center: https://support.tuya.com/en/help
-   Service & Support: https://service.console.tuya.com/8/2/list?source=content_feedback
