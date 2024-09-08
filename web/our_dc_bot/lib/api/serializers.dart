import 'package:built_collection/built_collection.dart';
import 'package:built_value/serializer.dart';
import 'package:built_value/standard_json_plugin.dart';

import 'package:our_dc_bot/api/exception.dart';
import 'model.dart';

part 'serializers.g.dart';

@SerializersFor([
  ErrorResponse,
  LoginCode,
  UserInformation,
  VerifyLoginCodeRequest,
  VerifyLoginCodeResponse,
  GetStickerByNameResponse,
  ListStickerResponse,
  Sticker,
  StickerImage,
  AddStickerRequest,
])
final Serializers serializers =
    (_$serializers.toBuilder()..addPlugin(StandardJsonPlugin())).build();
