// package: server
// file: proto/sara.proto

var proto_sara_pb = require("../proto/sara_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var pictureLinksService = (function () {
  function pictureLinksService() {}
  pictureLinksService.serviceName = "server.pictureLinksService";
  return pictureLinksService;
}());

pictureLinksService.pictureLinks = {
  methodName: "pictureLinks",
  service: pictureLinksService,
  requestStream: false,
  responseStream: false,
  requestType: proto_sara_pb.pictureLinksRequest,
  responseType: proto_sara_pb.pictureLinksResponse
};

exports.pictureLinksService = pictureLinksService;

function pictureLinksServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

pictureLinksServiceClient.prototype.pictureLinks = function pictureLinks(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(pictureLinksService.pictureLinks, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

exports.pictureLinksServiceClient = pictureLinksServiceClient;

