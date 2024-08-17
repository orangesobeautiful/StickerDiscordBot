// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'serializers.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

Serializers _$serializers = (new Serializers().toBuilder()
      ..add(AddStickerRequest.serializer)
      ..add(ErrorResponse.serializer)
      ..add(ListStickerResponse.serializer)
      ..add(LoginCode.serializer)
      ..add(Sticker.serializer)
      ..add(StickerImage.serializer)
      ..add(UserInformation.serializer)
      ..add(VerifyLoginCodeRequest.serializer)
      ..add(VerifyLoginCodeResponse.serializer)
      ..addBuilderFactory(
          const FullType(BuiltList, const [const FullType(Sticker)]),
          () => new ListBuilder<Sticker>())
      ..addBuilderFactory(
          const FullType(BuiltList, const [const FullType(StickerImage)]),
          () => new ListBuilder<StickerImage>())
      ..addBuilderFactory(
          const FullType(BuiltList, const [const FullType(String)]),
          () => new ListBuilder<String>()))
    .build();

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
