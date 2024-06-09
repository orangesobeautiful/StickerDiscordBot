import 'dart:convert';

import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

import 'package:our_dc_bot/api/serializers.dart';
import 'package:our_dc_bot/api/tojson.dart';

part 'model.g.dart';

abstract class LoginCode
    implements Built<LoginCode, LoginCodeBuilder>, ToJsoner {
  static Serializer<LoginCode> get serializer => _$loginCodeSerializer;

  @BuiltValueField(wireName: 'code')
  String get code;

  LoginCode._();

  factory LoginCode([void Function(LoginCodeBuilder) updates]) = _$LoginCode;

  @override
  String toJson() {
    return json.encode(serializers.serializeWith(LoginCode.serializer, this));
  }

  static LoginCode? fromJson(String jsonString) {
    return serializers.deserializeWith(
        LoginCode.serializer, json.decode(jsonString));
  }
}

abstract class UserInformation
    implements Built<UserInformation, UserInformationBuilder>, ToJsoner {
  static Serializer<UserInformation> get serializer =>
      _$userInformationSerializer;

  @BuiltValueField(wireName: 'id')
  int get id;

  @BuiltValueField(wireName: 'discord_id')
  String get discordId;

  @BuiltValueField(wireName: 'guild_id')
  String get guildId;

  @BuiltValueField(wireName: 'name')
  String get name;

  @BuiltValueField(wireName: 'avatar_url')
  String get avatarUrl;

  UserInformation._();

  factory UserInformation([void Function(UserInformationBuilder) updates]) =
      _$UserInformation;

  @override
  String toJson() {
    return json
        .encode(serializers.serializeWith(UserInformation.serializer, this));
  }

  static UserInformation? fromJson(String jsonString) {
    return serializers.deserializeWith(
        UserInformation.serializer, json.decode(jsonString));
  }
}

abstract class VerifyLoginCodeRequest
    implements
        Built<VerifyLoginCodeRequest, VerifyLoginCodeRequestBuilder>,
        ToJsoner {
  static Serializer<VerifyLoginCodeRequest> get serializer =>
      _$verifyLoginCodeRequestSerializer;

  @BuiltValueField(wireName: 'code')
  String get code;

  VerifyLoginCodeRequest._();

  factory VerifyLoginCodeRequest(
          [void Function(VerifyLoginCodeRequestBuilder) updates]) =
      _$VerifyLoginCodeRequest;

  @override
  String toJson() {
    return json.encode(
        serializers.serializeWith(VerifyLoginCodeRequest.serializer, this));
  }

  static VerifyLoginCodeRequest? fromJson(String jsonString) {
    return serializers.deserializeWith(
        VerifyLoginCodeRequest.serializer, json.decode(jsonString));
  }
}

abstract class VerifyLoginCodeResponse
    implements
        Built<VerifyLoginCodeResponse, VerifyLoginCodeResponseBuilder>,
        ToJsoner {
  static Serializer<VerifyLoginCodeResponse> get serializer =>
      _$verifyLoginCodeResponseSerializer;

  @BuiltValueField(wireName: 'is_verified')
  bool get isVerified;

  @BuiltValueField(wireName: 'token')
  String get token;

  VerifyLoginCodeResponse._();

  factory VerifyLoginCodeResponse(
          [void Function(VerifyLoginCodeResponseBuilder) updates]) =
      _$VerifyLoginCodeResponse;

  @override
  String toJson() {
    return json.encode(
        serializers.serializeWith(VerifyLoginCodeResponse.serializer, this));
  }

  static VerifyLoginCodeResponse? fromJson(String jsonString) {
    return serializers.deserializeWith(
        VerifyLoginCodeResponse.serializer, json.decode(jsonString));
  }
}
