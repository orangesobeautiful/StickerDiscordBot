// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'model.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

Serializer<LoginCode> _$loginCodeSerializer = new _$LoginCodeSerializer();
Serializer<UserInformation> _$userInformationSerializer =
    new _$UserInformationSerializer();
Serializer<VerifyLoginCodeRequest> _$verifyLoginCodeRequestSerializer =
    new _$VerifyLoginCodeRequestSerializer();
Serializer<VerifyLoginCodeResponse> _$verifyLoginCodeResponseSerializer =
    new _$VerifyLoginCodeResponseSerializer();
Serializer<GetStickerByNameResponse> _$getStickerByNameResponseSerializer =
    new _$GetStickerByNameResponseSerializer();
Serializer<ListStickerResponse> _$listStickerResponseSerializer =
    new _$ListStickerResponseSerializer();
Serializer<Sticker> _$stickerSerializer = new _$StickerSerializer();
Serializer<StickerImage> _$stickerImageSerializer =
    new _$StickerImageSerializer();
Serializer<AddStickerRequest> _$addStickerRequestSerializer =
    new _$AddStickerRequestSerializer();

class _$LoginCodeSerializer implements StructuredSerializer<LoginCode> {
  @override
  final Iterable<Type> types = const [LoginCode, _$LoginCode];
  @override
  final String wireName = 'LoginCode';

  @override
  Iterable<Object?> serialize(Serializers serializers, LoginCode object,
      {FullType specifiedType = FullType.unspecified}) {
    final result = <Object?>[
      'code',
      serializers.serialize(object.code, specifiedType: const FullType(String)),
    ];

    return result;
  }

  @override
  LoginCode deserialize(Serializers serializers, Iterable<Object?> serialized,
      {FullType specifiedType = FullType.unspecified}) {
    final result = new LoginCodeBuilder();

    final iterator = serialized.iterator;
    while (iterator.moveNext()) {
      final key = iterator.current! as String;
      iterator.moveNext();
      final Object? value = iterator.current;
      switch (key) {
        case 'code':
          result.code = serializers.deserialize(value,
              specifiedType: const FullType(String))! as String;
          break;
      }
    }

    return result.build();
  }
}

class _$UserInformationSerializer
    implements StructuredSerializer<UserInformation> {
  @override
  final Iterable<Type> types = const [UserInformation, _$UserInformation];
  @override
  final String wireName = 'UserInformation';

  @override
  Iterable<Object?> serialize(Serializers serializers, UserInformation object,
      {FullType specifiedType = FullType.unspecified}) {
    final result = <Object?>[
      'id',
      serializers.serialize(object.id, specifiedType: const FullType(int)),
      'discord_id',
      serializers.serialize(object.discordId,
          specifiedType: const FullType(String)),
      'guild_id',
      serializers.serialize(object.guildId,
          specifiedType: const FullType(String)),
      'name',
      serializers.serialize(object.name, specifiedType: const FullType(String)),
      'avatar_url',
      serializers.serialize(object.avatarUrl,
          specifiedType: const FullType(String)),
    ];

    return result;
  }

  @override
  UserInformation deserialize(
      Serializers serializers, Iterable<Object?> serialized,
      {FullType specifiedType = FullType.unspecified}) {
    final result = new UserInformationBuilder();

    final iterator = serialized.iterator;
    while (iterator.moveNext()) {
      final key = iterator.current! as String;
      iterator.moveNext();
      final Object? value = iterator.current;
      switch (key) {
        case 'id':
          result.id = serializers.deserialize(value,
              specifiedType: const FullType(int))! as int;
          break;
        case 'discord_id':
          result.discordId = serializers.deserialize(value,
              specifiedType: const FullType(String))! as String;
          break;
        case 'guild_id':
          result.guildId = serializers.deserialize(value,
              specifiedType: const FullType(String))! as String;
          break;
        case 'name':
          result.name = serializers.deserialize(value,
              specifiedType: const FullType(String))! as String;
          break;
        case 'avatar_url':
          result.avatarUrl = serializers.deserialize(value,
              specifiedType: const FullType(String))! as String;
          break;
      }
    }

    return result.build();
  }
}

class _$VerifyLoginCodeRequestSerializer
    implements StructuredSerializer<VerifyLoginCodeRequest> {
  @override
  final Iterable<Type> types = const [
    VerifyLoginCodeRequest,
    _$VerifyLoginCodeRequest
  ];
  @override
  final String wireName = 'VerifyLoginCodeRequest';

  @override
  Iterable<Object?> serialize(
      Serializers serializers, VerifyLoginCodeRequest object,
      {FullType specifiedType = FullType.unspecified}) {
    final result = <Object?>[
      'code',
      serializers.serialize(object.code, specifiedType: const FullType(String)),
    ];

    return result;
  }

  @override
  VerifyLoginCodeRequest deserialize(
      Serializers serializers, Iterable<Object?> serialized,
      {FullType specifiedType = FullType.unspecified}) {
    final result = new VerifyLoginCodeRequestBuilder();

    final iterator = serialized.iterator;
    while (iterator.moveNext()) {
      final key = iterator.current! as String;
      iterator.moveNext();
      final Object? value = iterator.current;
      switch (key) {
        case 'code':
          result.code = serializers.deserialize(value,
              specifiedType: const FullType(String))! as String;
          break;
      }
    }

    return result.build();
  }
}

class _$VerifyLoginCodeResponseSerializer
    implements StructuredSerializer<VerifyLoginCodeResponse> {
  @override
  final Iterable<Type> types = const [
    VerifyLoginCodeResponse,
    _$VerifyLoginCodeResponse
  ];
  @override
  final String wireName = 'VerifyLoginCodeResponse';

  @override
  Iterable<Object?> serialize(
      Serializers serializers, VerifyLoginCodeResponse object,
      {FullType specifiedType = FullType.unspecified}) {
    final result = <Object?>[
      'is_verified',
      serializers.serialize(object.isVerified,
          specifiedType: const FullType(bool)),
      'token',
      serializers.serialize(object.token,
          specifiedType: const FullType(String)),
    ];

    return result;
  }

  @override
  VerifyLoginCodeResponse deserialize(
      Serializers serializers, Iterable<Object?> serialized,
      {FullType specifiedType = FullType.unspecified}) {
    final result = new VerifyLoginCodeResponseBuilder();

    final iterator = serialized.iterator;
    while (iterator.moveNext()) {
      final key = iterator.current! as String;
      iterator.moveNext();
      final Object? value = iterator.current;
      switch (key) {
        case 'is_verified':
          result.isVerified = serializers.deserialize(value,
              specifiedType: const FullType(bool))! as bool;
          break;
        case 'token':
          result.token = serializers.deserialize(value,
              specifiedType: const FullType(String))! as String;
          break;
      }
    }

    return result.build();
  }
}

class _$GetStickerByNameResponseSerializer
    implements StructuredSerializer<GetStickerByNameResponse> {
  @override
  final Iterable<Type> types = const [
    GetStickerByNameResponse,
    _$GetStickerByNameResponse
  ];
  @override
  final String wireName = 'GetStickerByNameResponse';

  @override
  Iterable<Object?> serialize(
      Serializers serializers, GetStickerByNameResponse object,
      {FullType specifiedType = FullType.unspecified}) {
    final result = <Object?>[
      'sticker',
      serializers.serialize(object.sticker,
          specifiedType: const FullType(Sticker)),
    ];

    return result;
  }

  @override
  GetStickerByNameResponse deserialize(
      Serializers serializers, Iterable<Object?> serialized,
      {FullType specifiedType = FullType.unspecified}) {
    final result = new GetStickerByNameResponseBuilder();

    final iterator = serialized.iterator;
    while (iterator.moveNext()) {
      final key = iterator.current! as String;
      iterator.moveNext();
      final Object? value = iterator.current;
      switch (key) {
        case 'sticker':
          result.sticker.replace(serializers.deserialize(value,
              specifiedType: const FullType(Sticker))! as Sticker);
          break;
      }
    }

    return result.build();
  }
}

class _$ListStickerResponseSerializer
    implements StructuredSerializer<ListStickerResponse> {
  @override
  final Iterable<Type> types = const [
    ListStickerResponse,
    _$ListStickerResponse
  ];
  @override
  final String wireName = 'ListStickerResponse';

  @override
  Iterable<Object?> serialize(
      Serializers serializers, ListStickerResponse object,
      {FullType specifiedType = FullType.unspecified}) {
    final result = <Object?>[
      'total_count',
      serializers.serialize(object.totalCount,
          specifiedType: const FullType(int)),
      'stickers',
      serializers.serialize(object.stickers,
          specifiedType:
              const FullType(BuiltList, const [const FullType(Sticker)])),
    ];

    return result;
  }

  @override
  ListStickerResponse deserialize(
      Serializers serializers, Iterable<Object?> serialized,
      {FullType specifiedType = FullType.unspecified}) {
    final result = new ListStickerResponseBuilder();

    final iterator = serialized.iterator;
    while (iterator.moveNext()) {
      final key = iterator.current! as String;
      iterator.moveNext();
      final Object? value = iterator.current;
      switch (key) {
        case 'total_count':
          result.totalCount = serializers.deserialize(value,
              specifiedType: const FullType(int))! as int;
          break;
        case 'stickers':
          result.stickers.replace(serializers.deserialize(value,
                  specifiedType: const FullType(
                      BuiltList, const [const FullType(Sticker)]))!
              as BuiltList<Object?>);
          break;
      }
    }

    return result.build();
  }
}

class _$StickerSerializer implements StructuredSerializer<Sticker> {
  @override
  final Iterable<Type> types = const [Sticker, _$Sticker];
  @override
  final String wireName = 'Sticker';

  @override
  Iterable<Object?> serialize(Serializers serializers, Sticker object,
      {FullType specifiedType = FullType.unspecified}) {
    final result = <Object?>[
      'id',
      serializers.serialize(object.id, specifiedType: const FullType(int)),
      'sticker_name',
      serializers.serialize(object.stickerName,
          specifiedType: const FullType(String)),
      'images',
      serializers.serialize(object.images,
          specifiedType:
              const FullType(BuiltList, const [const FullType(StickerImage)])),
    ];

    return result;
  }

  @override
  Sticker deserialize(Serializers serializers, Iterable<Object?> serialized,
      {FullType specifiedType = FullType.unspecified}) {
    final result = new StickerBuilder();

    final iterator = serialized.iterator;
    while (iterator.moveNext()) {
      final key = iterator.current! as String;
      iterator.moveNext();
      final Object? value = iterator.current;
      switch (key) {
        case 'id':
          result.id = serializers.deserialize(value,
              specifiedType: const FullType(int))! as int;
          break;
        case 'sticker_name':
          result.stickerName = serializers.deserialize(value,
              specifiedType: const FullType(String))! as String;
          break;
        case 'images':
          result.images.replace(serializers.deserialize(value,
                  specifiedType: const FullType(
                      BuiltList, const [const FullType(StickerImage)]))!
              as BuiltList<Object?>);
          break;
      }
    }

    return result.build();
  }
}

class _$StickerImageSerializer implements StructuredSerializer<StickerImage> {
  @override
  final Iterable<Type> types = const [StickerImage, _$StickerImage];
  @override
  final String wireName = 'StickerImage';

  @override
  Iterable<Object?> serialize(Serializers serializers, StickerImage object,
      {FullType specifiedType = FullType.unspecified}) {
    final result = <Object?>[
      'id',
      serializers.serialize(object.id, specifiedType: const FullType(int)),
      'url',
      serializers.serialize(object.url, specifiedType: const FullType(String)),
    ];

    return result;
  }

  @override
  StickerImage deserialize(
      Serializers serializers, Iterable<Object?> serialized,
      {FullType specifiedType = FullType.unspecified}) {
    final result = new StickerImageBuilder();

    final iterator = serialized.iterator;
    while (iterator.moveNext()) {
      final key = iterator.current! as String;
      iterator.moveNext();
      final Object? value = iterator.current;
      switch (key) {
        case 'id':
          result.id = serializers.deserialize(value,
              specifiedType: const FullType(int))! as int;
          break;
        case 'url':
          result.url = serializers.deserialize(value,
              specifiedType: const FullType(String))! as String;
          break;
      }
    }

    return result.build();
  }
}

class _$AddStickerRequestSerializer
    implements StructuredSerializer<AddStickerRequest> {
  @override
  final Iterable<Type> types = const [AddStickerRequest, _$AddStickerRequest];
  @override
  final String wireName = 'AddStickerRequest';

  @override
  Iterable<Object?> serialize(Serializers serializers, AddStickerRequest object,
      {FullType specifiedType = FullType.unspecified}) {
    final result = <Object?>[
      'sticker_name',
      serializers.serialize(object.stickerName,
          specifiedType: const FullType(String)),
      'image_url',
      serializers.serialize(object.imageURL,
          specifiedType: const FullType(String)),
    ];

    return result;
  }

  @override
  AddStickerRequest deserialize(
      Serializers serializers, Iterable<Object?> serialized,
      {FullType specifiedType = FullType.unspecified}) {
    final result = new AddStickerRequestBuilder();

    final iterator = serialized.iterator;
    while (iterator.moveNext()) {
      final key = iterator.current! as String;
      iterator.moveNext();
      final Object? value = iterator.current;
      switch (key) {
        case 'sticker_name':
          result.stickerName = serializers.deserialize(value,
              specifiedType: const FullType(String))! as String;
          break;
        case 'image_url':
          result.imageURL = serializers.deserialize(value,
              specifiedType: const FullType(String))! as String;
          break;
      }
    }

    return result.build();
  }
}

class _$LoginCode extends LoginCode {
  @override
  final String code;

  factory _$LoginCode([void Function(LoginCodeBuilder)? updates]) =>
      (new LoginCodeBuilder()..update(updates))._build();

  _$LoginCode._({required this.code}) : super._() {
    BuiltValueNullFieldError.checkNotNull(code, r'LoginCode', 'code');
  }

  @override
  LoginCode rebuild(void Function(LoginCodeBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  LoginCodeBuilder toBuilder() => new LoginCodeBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is LoginCode && code == other.code;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, code.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'LoginCode')..add('code', code))
        .toString();
  }
}

class LoginCodeBuilder implements Builder<LoginCode, LoginCodeBuilder> {
  _$LoginCode? _$v;

  String? _code;
  String? get code => _$this._code;
  set code(String? code) => _$this._code = code;

  LoginCodeBuilder();

  LoginCodeBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _code = $v.code;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(LoginCode other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$LoginCode;
  }

  @override
  void update(void Function(LoginCodeBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  LoginCode build() => _build();

  _$LoginCode _build() {
    final _$result = _$v ??
        new _$LoginCode._(
            code: BuiltValueNullFieldError.checkNotNull(
                code, r'LoginCode', 'code'));
    replace(_$result);
    return _$result;
  }
}

class _$UserInformation extends UserInformation {
  @override
  final int id;
  @override
  final String discordId;
  @override
  final String guildId;
  @override
  final String name;
  @override
  final String avatarUrl;

  factory _$UserInformation([void Function(UserInformationBuilder)? updates]) =>
      (new UserInformationBuilder()..update(updates))._build();

  _$UserInformation._(
      {required this.id,
      required this.discordId,
      required this.guildId,
      required this.name,
      required this.avatarUrl})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(id, r'UserInformation', 'id');
    BuiltValueNullFieldError.checkNotNull(
        discordId, r'UserInformation', 'discordId');
    BuiltValueNullFieldError.checkNotNull(
        guildId, r'UserInformation', 'guildId');
    BuiltValueNullFieldError.checkNotNull(name, r'UserInformation', 'name');
    BuiltValueNullFieldError.checkNotNull(
        avatarUrl, r'UserInformation', 'avatarUrl');
  }

  @override
  UserInformation rebuild(void Function(UserInformationBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  UserInformationBuilder toBuilder() =>
      new UserInformationBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is UserInformation &&
        id == other.id &&
        discordId == other.discordId &&
        guildId == other.guildId &&
        name == other.name &&
        avatarUrl == other.avatarUrl;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, discordId.hashCode);
    _$hash = $jc(_$hash, guildId.hashCode);
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, avatarUrl.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'UserInformation')
          ..add('id', id)
          ..add('discordId', discordId)
          ..add('guildId', guildId)
          ..add('name', name)
          ..add('avatarUrl', avatarUrl))
        .toString();
  }
}

class UserInformationBuilder
    implements Builder<UserInformation, UserInformationBuilder> {
  _$UserInformation? _$v;

  int? _id;
  int? get id => _$this._id;
  set id(int? id) => _$this._id = id;

  String? _discordId;
  String? get discordId => _$this._discordId;
  set discordId(String? discordId) => _$this._discordId = discordId;

  String? _guildId;
  String? get guildId => _$this._guildId;
  set guildId(String? guildId) => _$this._guildId = guildId;

  String? _name;
  String? get name => _$this._name;
  set name(String? name) => _$this._name = name;

  String? _avatarUrl;
  String? get avatarUrl => _$this._avatarUrl;
  set avatarUrl(String? avatarUrl) => _$this._avatarUrl = avatarUrl;

  UserInformationBuilder();

  UserInformationBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _id = $v.id;
      _discordId = $v.discordId;
      _guildId = $v.guildId;
      _name = $v.name;
      _avatarUrl = $v.avatarUrl;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(UserInformation other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$UserInformation;
  }

  @override
  void update(void Function(UserInformationBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  UserInformation build() => _build();

  _$UserInformation _build() {
    final _$result = _$v ??
        new _$UserInformation._(
            id: BuiltValueNullFieldError.checkNotNull(
                id, r'UserInformation', 'id'),
            discordId: BuiltValueNullFieldError.checkNotNull(
                discordId, r'UserInformation', 'discordId'),
            guildId: BuiltValueNullFieldError.checkNotNull(
                guildId, r'UserInformation', 'guildId'),
            name: BuiltValueNullFieldError.checkNotNull(
                name, r'UserInformation', 'name'),
            avatarUrl: BuiltValueNullFieldError.checkNotNull(
                avatarUrl, r'UserInformation', 'avatarUrl'));
    replace(_$result);
    return _$result;
  }
}

class _$VerifyLoginCodeRequest extends VerifyLoginCodeRequest {
  @override
  final String code;

  factory _$VerifyLoginCodeRequest(
          [void Function(VerifyLoginCodeRequestBuilder)? updates]) =>
      (new VerifyLoginCodeRequestBuilder()..update(updates))._build();

  _$VerifyLoginCodeRequest._({required this.code}) : super._() {
    BuiltValueNullFieldError.checkNotNull(
        code, r'VerifyLoginCodeRequest', 'code');
  }

  @override
  VerifyLoginCodeRequest rebuild(
          void Function(VerifyLoginCodeRequestBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  VerifyLoginCodeRequestBuilder toBuilder() =>
      new VerifyLoginCodeRequestBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is VerifyLoginCodeRequest && code == other.code;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, code.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'VerifyLoginCodeRequest')
          ..add('code', code))
        .toString();
  }
}

class VerifyLoginCodeRequestBuilder
    implements Builder<VerifyLoginCodeRequest, VerifyLoginCodeRequestBuilder> {
  _$VerifyLoginCodeRequest? _$v;

  String? _code;
  String? get code => _$this._code;
  set code(String? code) => _$this._code = code;

  VerifyLoginCodeRequestBuilder();

  VerifyLoginCodeRequestBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _code = $v.code;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(VerifyLoginCodeRequest other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$VerifyLoginCodeRequest;
  }

  @override
  void update(void Function(VerifyLoginCodeRequestBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  VerifyLoginCodeRequest build() => _build();

  _$VerifyLoginCodeRequest _build() {
    final _$result = _$v ??
        new _$VerifyLoginCodeRequest._(
            code: BuiltValueNullFieldError.checkNotNull(
                code, r'VerifyLoginCodeRequest', 'code'));
    replace(_$result);
    return _$result;
  }
}

class _$VerifyLoginCodeResponse extends VerifyLoginCodeResponse {
  @override
  final bool isVerified;
  @override
  final String token;

  factory _$VerifyLoginCodeResponse(
          [void Function(VerifyLoginCodeResponseBuilder)? updates]) =>
      (new VerifyLoginCodeResponseBuilder()..update(updates))._build();

  _$VerifyLoginCodeResponse._({required this.isVerified, required this.token})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(
        isVerified, r'VerifyLoginCodeResponse', 'isVerified');
    BuiltValueNullFieldError.checkNotNull(
        token, r'VerifyLoginCodeResponse', 'token');
  }

  @override
  VerifyLoginCodeResponse rebuild(
          void Function(VerifyLoginCodeResponseBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  VerifyLoginCodeResponseBuilder toBuilder() =>
      new VerifyLoginCodeResponseBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is VerifyLoginCodeResponse &&
        isVerified == other.isVerified &&
        token == other.token;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, isVerified.hashCode);
    _$hash = $jc(_$hash, token.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'VerifyLoginCodeResponse')
          ..add('isVerified', isVerified)
          ..add('token', token))
        .toString();
  }
}

class VerifyLoginCodeResponseBuilder
    implements
        Builder<VerifyLoginCodeResponse, VerifyLoginCodeResponseBuilder> {
  _$VerifyLoginCodeResponse? _$v;

  bool? _isVerified;
  bool? get isVerified => _$this._isVerified;
  set isVerified(bool? isVerified) => _$this._isVerified = isVerified;

  String? _token;
  String? get token => _$this._token;
  set token(String? token) => _$this._token = token;

  VerifyLoginCodeResponseBuilder();

  VerifyLoginCodeResponseBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _isVerified = $v.isVerified;
      _token = $v.token;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(VerifyLoginCodeResponse other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$VerifyLoginCodeResponse;
  }

  @override
  void update(void Function(VerifyLoginCodeResponseBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  VerifyLoginCodeResponse build() => _build();

  _$VerifyLoginCodeResponse _build() {
    final _$result = _$v ??
        new _$VerifyLoginCodeResponse._(
            isVerified: BuiltValueNullFieldError.checkNotNull(
                isVerified, r'VerifyLoginCodeResponse', 'isVerified'),
            token: BuiltValueNullFieldError.checkNotNull(
                token, r'VerifyLoginCodeResponse', 'token'));
    replace(_$result);
    return _$result;
  }
}

class _$GetStickerByNameResponse extends GetStickerByNameResponse {
  @override
  final Sticker sticker;

  factory _$GetStickerByNameResponse(
          [void Function(GetStickerByNameResponseBuilder)? updates]) =>
      (new GetStickerByNameResponseBuilder()..update(updates))._build();

  _$GetStickerByNameResponse._({required this.sticker}) : super._() {
    BuiltValueNullFieldError.checkNotNull(
        sticker, r'GetStickerByNameResponse', 'sticker');
  }

  @override
  GetStickerByNameResponse rebuild(
          void Function(GetStickerByNameResponseBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  GetStickerByNameResponseBuilder toBuilder() =>
      new GetStickerByNameResponseBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is GetStickerByNameResponse && sticker == other.sticker;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, sticker.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'GetStickerByNameResponse')
          ..add('sticker', sticker))
        .toString();
  }
}

class GetStickerByNameResponseBuilder
    implements
        Builder<GetStickerByNameResponse, GetStickerByNameResponseBuilder> {
  _$GetStickerByNameResponse? _$v;

  StickerBuilder? _sticker;
  StickerBuilder get sticker => _$this._sticker ??= new StickerBuilder();
  set sticker(StickerBuilder? sticker) => _$this._sticker = sticker;

  GetStickerByNameResponseBuilder();

  GetStickerByNameResponseBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _sticker = $v.sticker.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(GetStickerByNameResponse other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$GetStickerByNameResponse;
  }

  @override
  void update(void Function(GetStickerByNameResponseBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  GetStickerByNameResponse build() => _build();

  _$GetStickerByNameResponse _build() {
    _$GetStickerByNameResponse _$result;
    try {
      _$result =
          _$v ?? new _$GetStickerByNameResponse._(sticker: sticker.build());
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'sticker';
        sticker.build();
      } catch (e) {
        throw new BuiltValueNestedFieldError(
            r'GetStickerByNameResponse', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

class _$ListStickerResponse extends ListStickerResponse {
  @override
  final int totalCount;
  @override
  final BuiltList<Sticker> stickers;

  factory _$ListStickerResponse(
          [void Function(ListStickerResponseBuilder)? updates]) =>
      (new ListStickerResponseBuilder()..update(updates))._build();

  _$ListStickerResponse._({required this.totalCount, required this.stickers})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(
        totalCount, r'ListStickerResponse', 'totalCount');
    BuiltValueNullFieldError.checkNotNull(
        stickers, r'ListStickerResponse', 'stickers');
  }

  @override
  ListStickerResponse rebuild(
          void Function(ListStickerResponseBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ListStickerResponseBuilder toBuilder() =>
      new ListStickerResponseBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ListStickerResponse &&
        totalCount == other.totalCount &&
        stickers == other.stickers;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, totalCount.hashCode);
    _$hash = $jc(_$hash, stickers.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ListStickerResponse')
          ..add('totalCount', totalCount)
          ..add('stickers', stickers))
        .toString();
  }
}

class ListStickerResponseBuilder
    implements Builder<ListStickerResponse, ListStickerResponseBuilder> {
  _$ListStickerResponse? _$v;

  int? _totalCount;
  int? get totalCount => _$this._totalCount;
  set totalCount(int? totalCount) => _$this._totalCount = totalCount;

  ListBuilder<Sticker>? _stickers;
  ListBuilder<Sticker> get stickers =>
      _$this._stickers ??= new ListBuilder<Sticker>();
  set stickers(ListBuilder<Sticker>? stickers) => _$this._stickers = stickers;

  ListStickerResponseBuilder();

  ListStickerResponseBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _totalCount = $v.totalCount;
      _stickers = $v.stickers.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ListStickerResponse other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$ListStickerResponse;
  }

  @override
  void update(void Function(ListStickerResponseBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ListStickerResponse build() => _build();

  _$ListStickerResponse _build() {
    _$ListStickerResponse _$result;
    try {
      _$result = _$v ??
          new _$ListStickerResponse._(
              totalCount: BuiltValueNullFieldError.checkNotNull(
                  totalCount, r'ListStickerResponse', 'totalCount'),
              stickers: stickers.build());
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'stickers';
        stickers.build();
      } catch (e) {
        throw new BuiltValueNestedFieldError(
            r'ListStickerResponse', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

class _$Sticker extends Sticker {
  @override
  final int id;
  @override
  final String stickerName;
  @override
  final BuiltList<StickerImage> images;

  factory _$Sticker([void Function(StickerBuilder)? updates]) =>
      (new StickerBuilder()..update(updates))._build();

  _$Sticker._(
      {required this.id, required this.stickerName, required this.images})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(id, r'Sticker', 'id');
    BuiltValueNullFieldError.checkNotNull(
        stickerName, r'Sticker', 'stickerName');
    BuiltValueNullFieldError.checkNotNull(images, r'Sticker', 'images');
  }

  @override
  Sticker rebuild(void Function(StickerBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  StickerBuilder toBuilder() => new StickerBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is Sticker &&
        id == other.id &&
        stickerName == other.stickerName &&
        images == other.images;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, stickerName.hashCode);
    _$hash = $jc(_$hash, images.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'Sticker')
          ..add('id', id)
          ..add('stickerName', stickerName)
          ..add('images', images))
        .toString();
  }
}

class StickerBuilder implements Builder<Sticker, StickerBuilder> {
  _$Sticker? _$v;

  int? _id;
  int? get id => _$this._id;
  set id(int? id) => _$this._id = id;

  String? _stickerName;
  String? get stickerName => _$this._stickerName;
  set stickerName(String? stickerName) => _$this._stickerName = stickerName;

  ListBuilder<StickerImage>? _images;
  ListBuilder<StickerImage> get images =>
      _$this._images ??= new ListBuilder<StickerImage>();
  set images(ListBuilder<StickerImage>? images) => _$this._images = images;

  StickerBuilder();

  StickerBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _id = $v.id;
      _stickerName = $v.stickerName;
      _images = $v.images.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(Sticker other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$Sticker;
  }

  @override
  void update(void Function(StickerBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  Sticker build() => _build();

  _$Sticker _build() {
    _$Sticker _$result;
    try {
      _$result = _$v ??
          new _$Sticker._(
              id: BuiltValueNullFieldError.checkNotNull(id, r'Sticker', 'id'),
              stickerName: BuiltValueNullFieldError.checkNotNull(
                  stickerName, r'Sticker', 'stickerName'),
              images: images.build());
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'images';
        images.build();
      } catch (e) {
        throw new BuiltValueNestedFieldError(
            r'Sticker', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

class _$StickerImage extends StickerImage {
  @override
  final int id;
  @override
  final String url;

  factory _$StickerImage([void Function(StickerImageBuilder)? updates]) =>
      (new StickerImageBuilder()..update(updates))._build();

  _$StickerImage._({required this.id, required this.url}) : super._() {
    BuiltValueNullFieldError.checkNotNull(id, r'StickerImage', 'id');
    BuiltValueNullFieldError.checkNotNull(url, r'StickerImage', 'url');
  }

  @override
  StickerImage rebuild(void Function(StickerImageBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  StickerImageBuilder toBuilder() => new StickerImageBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is StickerImage && id == other.id && url == other.url;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, url.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'StickerImage')
          ..add('id', id)
          ..add('url', url))
        .toString();
  }
}

class StickerImageBuilder
    implements Builder<StickerImage, StickerImageBuilder> {
  _$StickerImage? _$v;

  int? _id;
  int? get id => _$this._id;
  set id(int? id) => _$this._id = id;

  String? _url;
  String? get url => _$this._url;
  set url(String? url) => _$this._url = url;

  StickerImageBuilder();

  StickerImageBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _id = $v.id;
      _url = $v.url;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(StickerImage other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$StickerImage;
  }

  @override
  void update(void Function(StickerImageBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  StickerImage build() => _build();

  _$StickerImage _build() {
    final _$result = _$v ??
        new _$StickerImage._(
            id: BuiltValueNullFieldError.checkNotNull(
                id, r'StickerImage', 'id'),
            url: BuiltValueNullFieldError.checkNotNull(
                url, r'StickerImage', 'url'));
    replace(_$result);
    return _$result;
  }
}

class _$AddStickerRequest extends AddStickerRequest {
  @override
  final String stickerName;
  @override
  final String imageURL;

  factory _$AddStickerRequest(
          [void Function(AddStickerRequestBuilder)? updates]) =>
      (new AddStickerRequestBuilder()..update(updates))._build();

  _$AddStickerRequest._({required this.stickerName, required this.imageURL})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(
        stickerName, r'AddStickerRequest', 'stickerName');
    BuiltValueNullFieldError.checkNotNull(
        imageURL, r'AddStickerRequest', 'imageURL');
  }

  @override
  AddStickerRequest rebuild(void Function(AddStickerRequestBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  AddStickerRequestBuilder toBuilder() =>
      new AddStickerRequestBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is AddStickerRequest &&
        stickerName == other.stickerName &&
        imageURL == other.imageURL;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, stickerName.hashCode);
    _$hash = $jc(_$hash, imageURL.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'AddStickerRequest')
          ..add('stickerName', stickerName)
          ..add('imageURL', imageURL))
        .toString();
  }
}

class AddStickerRequestBuilder
    implements Builder<AddStickerRequest, AddStickerRequestBuilder> {
  _$AddStickerRequest? _$v;

  String? _stickerName;
  String? get stickerName => _$this._stickerName;
  set stickerName(String? stickerName) => _$this._stickerName = stickerName;

  String? _imageURL;
  String? get imageURL => _$this._imageURL;
  set imageURL(String? imageURL) => _$this._imageURL = imageURL;

  AddStickerRequestBuilder();

  AddStickerRequestBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _stickerName = $v.stickerName;
      _imageURL = $v.imageURL;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(AddStickerRequest other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$AddStickerRequest;
  }

  @override
  void update(void Function(AddStickerRequestBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  AddStickerRequest build() => _build();

  _$AddStickerRequest _build() {
    final _$result = _$v ??
        new _$AddStickerRequest._(
            stickerName: BuiltValueNullFieldError.checkNotNull(
                stickerName, r'AddStickerRequest', 'stickerName'),
            imageURL: BuiltValueNullFieldError.checkNotNull(
                imageURL, r'AddStickerRequest', 'imageURL'));
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
