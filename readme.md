### Struct

- Core

  - MessageWriter
  - MessageWriterPool
  - WriterCore：define by self, writer message to destination

- Custom

  - IWriterCore
  - IEncoder
  - IMessage：use to pack string and other field to be a message and put in writerCore
  - IMessageFactory：

### Design

#### log view

- string format

  ```
  [level] date time codedetail messagebody
  ```

- can set codedetail what content display

#### structure design

- use pool to decrease GC loading  ✅

- maybe use command pattern to refactor new message process ✅

- get log async and get log sync 

- how to let message factory can get some attribute that it is necessary to create new message.
  - Take a NewFactoryFunc into NewXLogFunc
  - execute it in the NewXLogFunc and store it into XLog

- link to another log api like loki promtail push api

### Loki promtail push pai

- [promtail push api](https://grafana.com/docs/loki/latest/community/design-documents/2020-02-promtail-push-api/)
