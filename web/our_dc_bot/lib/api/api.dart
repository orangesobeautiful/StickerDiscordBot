import 'dart:convert';

import 'package:flutter/material.dart';

import 'package:built_value/serializer.dart';
import 'package:compute/compute.dart';
import 'package:dio/browser.dart';
import 'package:dio/dio.dart';
import 'package:go_router/go_router.dart';

import 'package:our_dc_bot/api/exception.dart';
import 'package:our_dc_bot/api/model.dart';
import 'package:our_dc_bot/api/serializers.dart';
import 'package:our_dc_bot/api/tojson.dart';
import 'package:our_dc_bot/routers/enum.dart';

const _defaultConnectTimeout = Duration(seconds: 5);
const _defaultReceiveTimeout = Duration(seconds: 5);

typedef CallAPIErrorHandler = void Function(
    BuildContext context, dynamic error);
typedef UnauthorizedHandler = void Function(BuildContext context);
typedef OnErrorResponse = void Function(
    BuildContext context, int statusCode, ErrorResponse? errorResponse);

class API {
  API(
    Uri baseURL, {
    this.token = '',
    Duration connectTimeout = _defaultConnectTimeout,
    Duration receiveTimeout = _defaultReceiveTimeout,
  }) : dio = Dio(
          BaseOptions(
            baseUrl: baseURL.toString(),
            connectTimeout: connectTimeout,
            receiveTimeout: receiveTimeout,
            responseType: ResponseType.plain,
          ),
        )..httpClientAdapter = BrowserHttpClientAdapter(withCredentials: true) {
    dio.interceptors.add(InterceptorsWrapper(
      onRequest: (options, handler) => handler.next(options),
      onResponse: (response, handler) {
        return handler.next(response);
      },
    ));
  }

  String token;

  final Dio dio;

  Future<responseT> _normalJsonAPI<requestT, responseT>(
    String method,
    apiPath, {
    ToJsoner? requestData,
    Serializer<responseT>? respSerializer,
  }) async {
    Future<Response<dynamic>> apiCall =
        _httpSend(method, apiPath, requestData: requestData);

    // final Response response = await apiCall;
    Response response;
    try {
      response = await apiCall;
    } on DioException catch (dioException) {
      if (dioException.type == DioExceptionType.badResponse) {
        if (dioException.response!.statusCode == 401) {
          throw UnauthorizedException();
        }

        ErrorResponse? errorResponse;
        try {
          errorResponse = ErrorResponse.fromJson(
            dioException.response!.data,
          );
        } catch (e) {
          // ignore deserialize error
          print('Deserialize error: $e');
        }

        throw ErrorResponseException(
          dioException.response!.statusCode!,
          errorResponse,
        );
      }

      throw CallAPIException(dioException);
    } catch (e) {
      throw CallAPIException(e);
    }

    if (respSerializer == null) {
      return null as responseT;
    }

    responseT? data;
    try {
      data = await _isolateDeserializeJson(respSerializer, response.data);
    } catch (e) {
      throw FormatResponseDataException(e);
    }

    return data!;
  }

  Future<Response<dioResponseT>> _httpSend<requestT, dioResponseT>(
    String method,
    apiPath, {
    // ignore: avoid_unused_constructor_parameters
    ToJsoner? requestData,
  }) {
    Map<String, String> headers = {};
    if (token.isNotEmpty) {
      headers['Authorization'] = 'Bearer $token';
    }

    Object? body;
    if (method != 'GET' && requestData != null) {
      body = requestData.toJson();
    }

    switch (method) {
      case 'GET':
        return dio.get(apiPath, options: Options(headers: headers));
      case 'POST':
        return dio.post(apiPath,
            data: body,
            options: Options(
              headers: headers,
            ));
      case 'PUT':
        return dio.put(apiPath,
            data: body,
            options: Options(
              headers: headers,
            ));
      case 'PATCH':
        return dio.patch(apiPath,
            data: body,
            options: Options(
              headers: headers,
            ));
      case 'DELETE':
        return dio.delete(apiPath,
            data: body,
            options: Options(
              headers: headers,
            ));
      default:
        throw UnsupportedError('Unsupported method: $method');
    }
  }

  Future<LoginCode> getLoginCode() async {
    return _normalJsonAPI(
      'GET',
      '/api/v1/login-code',
      respSerializer: LoginCode.serializer,
    );
  }

  Future<UserInformation> getSelfInformation() async {
    return _normalJsonAPI(
      'GET',
      '/api/v1/me',
      respSerializer: UserInformation.serializer,
    );
  }

  Future<VerifyLoginCodeResponse> verifyLoginCode(String code) async {
    VerifyLoginCodeResponse verifyResult = await _normalJsonAPI(
      'POST',
      '/api/v1/login-code',
      requestData: VerifyLoginCodeRequest((b) => b..code = code),
      respSerializer: VerifyLoginCodeResponse.serializer,
    );

    if (verifyResult.isVerified) {
      token = verifyResult.token;
    }

    return verifyResult;
  }
}

Future<T> _isolateDeserializeJson<T>(
  Serializer<T> serializer,
  String responseData,
) async {
  return compute(_deserializeJson, (
    serializer: serializer,
    responseData: responseData,
  ));
}

T _deserializeJson<T>(({Serializer<T> serializer, String responseData}) arg) {
  T? data;

  try {
    data = serializers.deserializeWith(
        arg.serializer, jsonDecode(arg.responseData));
  } catch (e) {
    throw FormatResponseDataException(e);
  }

  return data as T;
}

CallAPIErrorHandler _defaultGlobalOnCallAPIError =
    (BuildContext context, dynamic error) {
  _defaultErrorSnackbar(context, 'Call API Error: $error');
};

UnauthorizedHandler _defaultGlobalOnUnauthorized = (BuildContext context) {
  context.goNamed(RouterName.signIn.name);
};

OnErrorResponse _defaultOnErrorResponse =
    (BuildContext context, int statusCode, ErrorResponse? errorResponse) {
  _defaultErrorSnackbar(
      context, '$statusCode: ${errorResponse?.message ?? 'Unknown Error'}');
};

void _defaultErrorSnackbar(BuildContext context, String message) {
  ScaffoldMessenger.of(context).showSnackBar(
    SnackBar(
      content: Text('Error: $message'),
      backgroundColor: Colors.red,
      duration: const Duration(seconds: 3),
      behavior: SnackBarBehavior.floating,
    ),
  );
}

class UIAPIHandler {
  UIAPIHandler(
    Uri baseURL, {
    String token = '',
    Duration connectTimeout = _defaultConnectTimeout,
    Duration receiveTimeout = _defaultReceiveTimeout,
    CallAPIErrorHandler? globalOnCallAPIError,
    UnauthorizedHandler? globalOnUnauthorized,
    OnErrorResponse? globalOnErrorResponse,
  })  : api = API(
          baseURL,
          token: token,
          connectTimeout: connectTimeout,
          receiveTimeout: receiveTimeout,
        ),
        _globalOnCallAPIError =
            globalOnCallAPIError ?? _defaultGlobalOnCallAPIError,
        _globalOnUnauthorized =
            globalOnUnauthorized ?? _defaultGlobalOnUnauthorized,
        _globalOnErrorResponse =
            globalOnErrorResponse ?? _defaultOnErrorResponse;

  final API api;

  final CallAPIErrorHandler _globalOnCallAPIError;
  final UnauthorizedHandler _globalOnUnauthorized;
  final OnErrorResponse _globalOnErrorResponse;

  Future<T> call<T>(
    BuildContext context,
    Future<T> Function(API api) apiCall,
  ) async {
    try {
      return await apiCall(api);
    } on ErrorResponseException catch (e) {
      if (context.mounted) {
        _globalOnErrorResponse(context, e.statusCode, e.errorResponse);
      }

      rethrow;
    } on UnauthorizedException {
      if (context.mounted) {
        _globalOnUnauthorized(context);
      }

      rethrow;
    } on FormatResponseDataException catch (e) {
      if (context.mounted) {
        _globalOnCallAPIError(context, e);
      }

      rethrow;
    } on CallAPIException catch (e) {
      if (context.mounted) {
        _globalOnCallAPIError(context, e);
      }

      rethrow;
    }
  }
}
