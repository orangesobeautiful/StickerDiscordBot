import 'package:flutter/material.dart';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_web_plugins/url_strategy.dart';

import 'package:our_dc_bot/routers/router.dart';

void main() {
  usePathUrlStrategy();

  runApp(
    const ProviderScope(
      child: MyApp(),
    ),
  );
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp.router(
      title: 'Our Discord Bot',
      theme: ThemeData(
        brightness: Brightness.dark,
      ),
      routerConfig: newRouter(),
    );
  }
}
