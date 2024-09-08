import 'dart:convert';

import 'package:built_value/built_value.dart';
import 'package:built_collection/built_collection.dart';
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

abstract class GetStickerByNameResponse
    implements
        Built<GetStickerByNameResponse, GetStickerByNameResponseBuilder>,
        ToJsoner {
  static Serializer<GetStickerByNameResponse> get serializer =>
      _$getStickerByNameResponseSerializer;

  @BuiltValueField(wireName: 'sticker')
  Sticker get sticker;

  GetStickerByNameResponse._();

  factory GetStickerByNameResponse(
          [void Function(GetStickerByNameResponseBuilder) updates]) =
      _$GetStickerByNameResponse;

  @override
  String toJson() {
    return json.encode(
        serializers.serializeWith(GetStickerByNameResponse.serializer, this));
  }

  static GetStickerByNameResponse? fromJson(String jsonString) {
    return serializers.deserializeWith(
        GetStickerByNameResponse.serializer, json.decode(jsonString));
  }
}

abstract class ListStickerResponse
    implements
        Built<ListStickerResponse, ListStickerResponseBuilder>,
        ToJsoner {
  static Serializer<ListStickerResponse> get serializer =>
      _$listStickerResponseSerializer;

  @BuiltValueField(wireName: 'total_count')
  int get totalCount;

  @BuiltValueField(wireName: 'stickers')
  BuiltList<Sticker> get stickers;

  ListStickerResponse._();

  factory ListStickerResponse(
          [void Function(ListStickerResponseBuilder) updates]) =
      _$ListStickerResponse;

  @override
  String toJson() {
    return json.encode(
        serializers.serializeWith(ListStickerResponse.serializer, this));
  }

  static ListStickerResponse? fromJson(String jsonString) {
    return serializers.deserializeWith(
        ListStickerResponse.serializer, json.decode(jsonString));
  }
}

abstract class Sticker implements Built<Sticker, StickerBuilder>, ToJsoner {
  static Serializer<Sticker> get serializer => _$stickerSerializer;

  @BuiltValueField(wireName: 'id')
  int get id;

  @BuiltValueField(wireName: 'sticker_name')
  String get stickerName;

  @BuiltValueField(wireName: 'images')
  BuiltList<StickerImage> get images;

  Sticker._();

  factory Sticker([void Function(StickerBuilder) updates]) = _$Sticker;

  @override
  String toJson() {
    return json.encode(serializers.serializeWith(Sticker.serializer, this));
  }

  static Sticker? fromJson(String jsonString) {
    return serializers.deserializeWith(
        Sticker.serializer, json.decode(jsonString));
  }
}

abstract class StickerImage
    implements Built<StickerImage, StickerImageBuilder>, ToJsoner {
  static Serializer<StickerImage> get serializer => _$stickerImageSerializer;

  @BuiltValueField(wireName: 'id')
  int get id;

  @BuiltValueField(wireName: 'url')
  String get url;

  StickerImage._();

  factory StickerImage([void Function(StickerImageBuilder) updates]) =
      _$StickerImage;

  @override
  String toJson() {
    return json
        .encode(serializers.serializeWith(StickerImage.serializer, this));
  }

  static StickerImage? fromJson(String jsonString) {
    return serializers.deserializeWith(
        StickerImage.serializer, json.decode(jsonString));
  }
}

abstract class AddStickerRequest
    implements Built<AddStickerRequest, AddStickerRequestBuilder>, ToJsoner {
  static Serializer<AddStickerRequest> get serializer =>
      _$addStickerRequestSerializer;

  @BuiltValueField(wireName: 'sticker_name')
  String get stickerName;

  @BuiltValueField(wireName: 'image_url')
  String get imageURL;

  AddStickerRequest._();

  factory AddStickerRequest([void Function(AddStickerRequestBuilder) updates]) =
      _$AddStickerRequest;

  @override
  String toJson() {
    return json
        .encode(serializers.serializeWith(AddStickerRequest.serializer, this));
  }

  static AddStickerRequest? fromJson(String jsonString) {
    return serializers.deserializeWith(
        AddStickerRequest.serializer, json.decode(jsonString));
  }
}
