### Struct

- XMessage：use to pack string and other field to be a message and put in writerCore
- WriterCore：define by self, writer message to destination
- XEncoder
- XMessageWriter
- XMessageWriterPool

### Design

- string format

  ```
  [level] date time codedetail messagebody
  ```

- use pool to decrease GC loading

- maybe use command pattern to refactor new message process

- how to let message factory can get some attribute that it is necessary to create new message.
  - Take a NewFactoryFunc into NewXLogFunc
  - execute it in the NewXLogFunc and store it into XLog

### Loki promtail push pai

- [promtail push api](https://grafana.com/docs/loki/latest/community/design-documents/2020-02-promtail-push-api/)
