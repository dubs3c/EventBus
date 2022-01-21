# EventBus

Tiny library for defining topics which can be subscribed to using channels. Data sent to a specific topic will be received by all functions that subscribed to the topic. EventBus only sets up the infrastructure for sending data to topics, it's up to the user to define how the functions should receive and process the data.

## Examples

Check [bus_test.go](bus_test.go) and [examples](examples/) to see how the library can be used.

## Contributing
1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :D