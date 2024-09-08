import 'dart:convert';

import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

import 'package:our_dc_bot/api/serializers.dart';
import 'package:our_dc_bot/api/tojson.dart';

part 'exception.g.dart';

class UnauthorizedException implements Exception {}

class NotFoundException implements Exception {}

class ErrorResponseException implements Exception {
  final int statusCode;
  final ErrorResponse? errorResponse;

  ErrorResponseException(this.statusCode, this.errorResponse);
}

class FormatResponseDataException implements Exception {
  final dynamic source;

  FormatResponseDataException(this.source);

  @override
  String toString() {
    return 'ResponseDataFormatException: $source';
  }
}

class CallAPIException implements Exception {
  final dynamic source;

  CallAPIException(this.source);

  @override
  String toString() {
    return 'CallAPIException: $source';
  }
}

abstract class ErrorResponse
    implements Built<ErrorResponse, ErrorResponseBuilder>, ToJsoner {
  static Serializer<ErrorResponse> get serializer => _$errorResponseSerializer;

  @BuiltValueField(wireName: 'Message')
  String get message;

  @BuiltValueField(wireName: 'Details')
  BuiltList<String> get details;

  ErrorResponse._();

  factory ErrorResponse([void Function(ErrorResponseBuilder) updates]) =
      _$ErrorResponse;

  @override
  String toJson() {
    return json.encoder
        .convert(serializers.serializeWith(ErrorResponse.serializer, this));
  }

  static ErrorResponse? fromJson(String jsonString) {
    return serializers.deserializeWith(
        ErrorResponse.serializer, json.decode(jsonString));
  }
}
