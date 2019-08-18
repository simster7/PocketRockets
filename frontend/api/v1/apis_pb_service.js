// package: v1
// file: api/v1/apis.proto

var api_v1_apis_pb = require("../../api/v1/apis_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var PokerService = (function () {
  function PokerService() {}
  PokerService.serviceName = "v1.PokerService";
  return PokerService;
}());

PokerService.StartGame = {
  methodName: "StartGame",
  service: PokerService,
  requestStream: false,
  responseStream: false,
  requestType: api_v1_apis_pb.StartGameRequest,
  responseType: api_v1_apis_pb.OperationResponse
};

PokerService.AddPlayer = {
  methodName: "AddPlayer",
  service: PokerService,
  requestStream: false,
  responseStream: false,
  requestType: api_v1_apis_pb.AddPlayerRequest,
  responseType: api_v1_apis_pb.OperationResponse
};

PokerService.SitPlayer = {
  methodName: "SitPlayer",
  service: PokerService,
  requestStream: false,
  responseStream: false,
  requestType: api_v1_apis_pb.SitPlayerRequest,
  responseType: api_v1_apis_pb.OperationResponse
};

PokerService.StandPlayer = {
  methodName: "StandPlayer",
  service: PokerService,
  requestStream: false,
  responseStream: false,
  requestType: api_v1_apis_pb.StandPlayerRequest,
  responseType: api_v1_apis_pb.OperationResponse
};

PokerService.TakeAction = {
  methodName: "TakeAction",
  service: PokerService,
  requestStream: false,
  responseStream: false,
  requestType: api_v1_apis_pb.TakeActionRequest,
  responseType: api_v1_apis_pb.OperationResponse
};

PokerService.DealHand = {
  methodName: "DealHand",
  service: PokerService,
  requestStream: false,
  responseStream: false,
  requestType: api_v1_apis_pb.DealHandRequest,
  responseType: api_v1_apis_pb.OperationResponse
};

PokerService.GetPlayerState = {
  methodName: "GetPlayerState",
  service: PokerService,
  requestStream: false,
  responseStream: false,
  requestType: api_v1_apis_pb.GetPlayerStateRequest,
  responseType: api_v1_apis_pb.PlayerState
};

exports.PokerService = PokerService;

function PokerServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

PokerServiceClient.prototype.startGame = function startGame(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(PokerService.StartGame, {
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

PokerServiceClient.prototype.addPlayer = function addPlayer(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(PokerService.AddPlayer, {
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

PokerServiceClient.prototype.sitPlayer = function sitPlayer(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(PokerService.SitPlayer, {
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

PokerServiceClient.prototype.standPlayer = function standPlayer(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(PokerService.StandPlayer, {
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

PokerServiceClient.prototype.takeAction = function takeAction(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(PokerService.TakeAction, {
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

PokerServiceClient.prototype.dealHand = function dealHand(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(PokerService.DealHand, {
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

PokerServiceClient.prototype.getPlayerState = function getPlayerState(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(PokerService.GetPlayerState, {
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

exports.PokerServiceClient = PokerServiceClient;

