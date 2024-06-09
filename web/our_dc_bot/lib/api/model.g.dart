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

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
