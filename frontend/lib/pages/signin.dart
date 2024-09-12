import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:intl/intl.dart';

import 'package:our_dc_bot/api/api.dart';
import 'package:our_dc_bot/global/global.dart';
import 'package:our_dc_bot/routers/enum.dart';

const Duration _checkLoginCodePeriod = Duration(seconds: 5);

class SignInPage extends ConsumerStatefulWidget {
  const SignInPage({super.key});

  @override
  SignInPageState createState() => SignInPageState();
}

class SignInPageState extends ConsumerState<SignInPage> {
  String _verifyCode = '';
  var effectiveSeconds = 90;

  Timer _checkLoginCodeTimer = Timer(Duration.zero, () {});

  @override
  void initState() {
    super.initState();

    _checkLoginCodeTimer = Timer.periodic(
      _checkLoginCodePeriod,
      (timer) {
        ref
            .read<UIAPIHandler>(apiHandlerProvider)
            .call(context, (api) => api.verifyLoginCode(_verifyCode))
            .then((verifyResult) {
          if (verifyResult.isVerified) {
            timer.cancel();
            context.goNamed(RouterName.myInfo.name);
          }
        });
      },
    );
  }

  @override
  void dispose() {
    _checkLoginCodeTimer.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    AuthState authState = ref.read(authStateNotifierProvider);
    if (authState.isLogin) {
      context.goNamed(RouterName.myInfo.name);
    }

    if (_verifyCode.isEmpty) {
      ref
          .read<UIAPIHandler>(apiHandlerProvider)
          .call(context, (api) => api.getLoginCode())
          .then(
        (value) {
          setState(() {
            _verifyCode = value.code;
          });
        },
      );
    }

    return Scaffold(
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: <Widget>[
            Text(
              '請在要登入的伺服器輸入以下驗證指令',
              style: Theme.of(context).textTheme.titleLarge,
            ),
            const SizedBox(height: 20),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Text('/web-login code:$_verifyCode'),
                const SizedBox(width: 10),
                _CountDownCopyButton(
                  copyContent: '/web-login code:$_verifyCode',
                  initEffectTotalSeconds: effectiveSeconds,
                  onRefreshPressed: () {
                    ref
                        .read<UIAPIHandler>(apiHandlerProvider)
                        .call(context, (api) => api.getLoginCode())
                        .then(
                      (value) {
                        setState(() {
                          _verifyCode = value.code;
                        });
                      },
                    );
                  },
                ),
              ],
            ),
            const SizedBox(height: 20),
          ],
        ),
      ),
    );
  }
}

class _CountDownCopyButton extends StatefulWidget {
  const _CountDownCopyButton(
      {required this.copyContent,
      required this.initEffectTotalSeconds,
      required this.onRefreshPressed});

  final String copyContent;
  final int initEffectTotalSeconds;
  final VoidCallback onRefreshPressed;

  @override
  State<_CountDownCopyButton> createState() => _CountDownCopyButtonState();
}

class _CountDownCopyButtonState extends State<_CountDownCopyButton> {
  int _effectTotalSeconds = 0;

  bool _isCountingDown = true;
  Timer _countDownTimer = Timer(Duration.zero, () {});

  @override
  void initState() {
    super.initState();

    _effectTotalSeconds = widget.initEffectTotalSeconds;

    _setEffectTotalSeconds(_effectTotalSeconds);
  }

  void _setEffectTotalSeconds(int seconds) {
    _effectTotalSeconds = seconds;
    _isCountingDown = true;
    caculateEffectTime();
    _newCowntDownTimer();
    setState(() {});
  }

  void _newCowntDownTimer() {
    _countDownTimer.cancel();
    _countDownTimer = Timer.periodic(
      const Duration(seconds: 1),
      (timer) {
        setState(() {
          _effectTotalSeconds -= 1;
          caculateEffectTime();

          if (_effectTotalSeconds <= 0) {
            _isCountingDown = false;
            setState(() {});
            timer.cancel();
          }
        });
      },
    );
  }

  @override
  void dispose() {
    _countDownTimer.cancel();
    super.dispose();
  }

  String _showEffectMinutes = '00';
  String _showEffectSeconds = '00';

  void caculateEffectTime() {
    _showEffectMinutes = NumberFormat('00').format(_effectTotalSeconds ~/ 60);
    _showEffectSeconds = NumberFormat('00').format(_effectTotalSeconds % 60);
  }

  @override
  Widget build(BuildContext context) {
    if (!_isCountingDown) {
      return ElevatedButton(
        child: const Icon(Icons.refresh),
        onPressed: () {
          widget.onRefreshPressed();
          _setEffectTotalSeconds(widget.initEffectTotalSeconds);
        },
      );
    }

    return ElevatedButton(
      child: Row(
        children: [
          const Icon(Icons.copy),
          const SizedBox(width: 10),
          Text(
            '$_showEffectMinutes:$_showEffectSeconds 後失效',
          ),
        ],
      ),
      onPressed: () {
        Clipboard.setData(
          ClipboardData(text: widget.copyContent),
        );
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            width: 200,
            content: Text(
              '複製成功',
              textAlign: TextAlign.center,
            ),
            backgroundColor: Colors.green,
            behavior: SnackBarBehavior.floating,
            duration: Duration(seconds: 1, milliseconds: 500),
          ),
        );
      },
    );
  }
}
