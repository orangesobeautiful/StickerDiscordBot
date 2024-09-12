import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter/gestures.dart';
import 'package:flutter_staggered_grid_view/flutter_staggered_grid_view.dart';
import 'package:number_paginator/number_paginator.dart';
import 'package:our_dc_bot/api/exception.dart';

import 'package:our_dc_bot/global/global.dart';
import 'package:our_dc_bot/api/api.dart';
import 'package:our_dc_bot/api/model.dart';

class StickerManagerPage extends ConsumerStatefulWidget {
  const StickerManagerPage({super.key});

  @override
  StickerManagerPageState createState() => StickerManagerPageState();
}

class StickerManagerPageState extends ConsumerState<StickerManagerPage>
    with TickerProviderStateMixin {
  final int pageSize = 5;

  String _searchText = '';
  int _currentPage = 1;
  int _numPages = 1;

  @override
  void initState() {
    super.initState();
  }

  @override
  void dispose() {
    super.dispose();
  }

  void _onSearchTextChanged(String text) {
    setState(() {
      _searchText = text;
      _currentPage = 1;
    });
  }

  @override
  Widget build(BuildContext context) {
    final userInfo = ref.read(authStateNotifierProvider).user;
    final api = ref.read<UIAPIHandler>(apiHandlerProvider);

    return Column(
      children: <Widget>[
        Padding(
          padding: const EdgeInsets.only(top: 32.0, bottom: 48.0),
          child: StickerManagerPageUpper(
            onSearchTextChanged: _onSearchTextChanged,
          ),
        ),
        Expanded(
          child: FutureBuilder(
            future: api.call(
              context,
              (api) => api.listSticker(
                userInfo.guildId,
                _currentPage,
                pageSize,
                search: _searchText,
              ),
            ),
            builder: (context, snapshot) {
              if (snapshot.connectionState == ConnectionState.waiting) {
                return const Center(
                  child: CircularProgressIndicator(),
                );
              }

              if (snapshot.hasError) {
                return const Center(
                  child: Text('Failed to load stickers'),
                );
              }

              final listStickerResp = snapshot.data as ListStickerResponse;

              _numPages = 1;
              if (listStickerResp.totalCount != 0) {
                _numPages = (listStickerResp.totalCount / pageSize).ceil();
              }

              return Column(
                children: <Widget>[
                  Expanded(
                    child: Center(
                      child: StickersView(
                        stickers: listStickerResp.stickers.toList(),
                        toRefresh: () {
                          setState(() {});
                        },
                      ),
                    ),
                  ),
                  SizedBox(
                    width: 500,
                    child: NumberPaginator(
                      numberPages: _numPages,
                      initialPage: _currentPage - 1,
                      onPageChange: (int index) {
                        setState(() {
                          _currentPage = index + 1;
                        });
                      },
                    ),
                  ),
                ],
              );
            },
          ),
        ),
      ],
    );
  }
}

class StickerManagerPageUpper extends ConsumerStatefulWidget {
  const StickerManagerPageUpper({
    super.key,
    required this.onSearchTextChanged,
  });

  final Function(String) onSearchTextChanged;

  @override
  StickerManagerPageUpperState createState() => StickerManagerPageUpperState();
}

class StickerManagerPageUpperState
    extends ConsumerState<StickerManagerPageUpper> {
  final _searchController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: <Widget>[
        const SizedBox(),
        SearchBar(
          controller: _searchController,
          padding: const WidgetStatePropertyAll<EdgeInsets>(
            EdgeInsets.symmetric(horizontal: 16.0),
          ),
          onSubmitted: (value) {
            widget.onSearchTextChanged(value);
          },
          leading: const Icon(Icons.search),
        ),
        FilledButton.tonalIcon(
          onPressed: () {
            showStickerManagaerDialog(context, '').then((value) {
              if (value != null) {
                setState(() {});
              }
            });
          },
          icon: const Icon(Icons.add),
          label: const Text('新增貼圖'),
          style: ElevatedButton.styleFrom(),
        ),
      ],
    );
  }
}

Future<bool?> showStickerManagaerDialog(BuildContext context, String name,
    {bool edit = false, int stickerID = 0}) {
  return showDialog<bool>(
    context: context,
    builder: (BuildContext context) {
      return AddStickerDialog(initName: name, edit: edit, stickerID: stickerID);
    },
  );
}

class AddStickerDialog extends ConsumerStatefulWidget {
  const AddStickerDialog(
      {super.key, String? initName, this.edit = false, this.stickerID = 0})
      : initName = initName ?? '';

  final String initName;
  final bool edit;
  final int stickerID;

  @override
  AddStickerDialogState createState() => AddStickerDialogState();
}

class AddStickerDialogState extends ConsumerState<AddStickerDialog> {
  final _stickerNameController = TextEditingController();

  @override
  void initState() {
    super.initState();

    _stickerNameController.text = widget.initName;
  }

  @override
  Widget build(BuildContext context) {
    Widget stickerManageBlock = Expanded(
      child: Center(
        child: Text(
          '輸入貼圖名稱',
          style: Theme.of(context).textTheme.titleLarge,
        ),
      ),
    );

    if (_stickerNameController.text.isNotEmpty) {
      stickerManageBlock = StickerDialogManageBlock(
        stickerName: _stickerNameController.text,
      );
    }

    return Dialog(
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Stack(
          children: <Widget>[
            Positioned(
              right: 0,
              top: 0,
              child: IconButton(
                icon: const Icon(Icons.delete),
                color: Colors.red,
                onPressed: () {
                  showDialog<bool>(
                    context: context,
                    builder: (BuildContext context) {
                      return AlertDialog(
                        title: const Text('確定要刪除嗎？'),
                        actions: <Widget>[
                          TextButton(
                            onPressed: () {
                              Navigator.of(context).pop(false);
                            },
                            child: const Text('取消'),
                          ),
                          TextButton(
                            onPressed: () async {
                              final api =
                                  ref.read<UIAPIHandler>(apiHandlerProvider);

                              await api.call(context, (api) {
                                return api.deleteSticker(widget.stickerID);
                              });

                              if (!context.mounted) {
                                return;
                              }

                              Navigator.of(context).pop(true);
                            },
                            child: const Text('刪除'),
                          ),
                        ],
                      );
                    },
                  ).then((bool? deleted) {
                    if (deleted == true && context.mounted) {
                      Navigator.of(context).pop(true);
                    }
                  });
                },
              ),
            ),
            Column(
              children: <Widget>[
                Text('貼圖管理', style: Theme.of(context).textTheme.titleLarge),
                const SizedBox(height: 50),
                SizedBox(
                  height: 50,
                  child: widget.edit
                      ? Text(widget.initName,
                          style: Theme.of(context).textTheme.titleLarge)
                      : TextField(
                          controller: _stickerNameController,
                          decoration: InputDecoration(
                            labelText: '貼圖名稱',
                            alignLabelWithHint: true,
                            enabled: !widget.edit,
                          ),
                          onChanged: (String value) {
                            setState(() {});
                          },
                        ),
                ),
                const SizedBox(height: 50),
                stickerManageBlock,
              ],
            ),
          ],
        ),
      ),
    );
  }
}

class StickerDialogManageBlock extends ConsumerStatefulWidget {
  const StickerDialogManageBlock({
    super.key,
    required this.stickerName,
  });

  final String stickerName;

  @override
  StickerDialogManageBlockState createState() =>
      StickerDialogManageBlockState();
}

class StickerDialogManageBlockState
    extends ConsumerState<StickerDialogManageBlock> {
  final _imageURLController = TextEditingController();

  Future<List<StickerImage>> getStickerImages(
      BuildContext context, String name) async {
    final userInfo = ref.read(authStateNotifierProvider).user;

    GetStickerByNameResponse response;

    try {
      response = await ref
          .read<UIAPIHandler>(apiHandlerProvider)
          .call(context, (api) => api.getStickerByName(userInfo.guildId, name));
    } on NotFoundException {
      return <StickerImage>[];
    }

    return response.sticker.images.toList();
  }

  Future<void> addImage(String imageURL) async {
    final userInfo = ref.read(authStateNotifierProvider).user;

    await ref.read<UIAPIHandler>(apiHandlerProvider).call(
          context,
          (api) =>
              api.addSticker(userInfo.guildId, widget.stickerName, imageURL),
        );

    if (!mounted) {
      return;
    }
  }

  @override
  Widget build(BuildContext context) {
    return FutureBuilder(
      future: getStickerImages(context, widget.stickerName),
      builder:
          (BuildContext context, AsyncSnapshot<List<StickerImage>> snapshot) {
        if (snapshot.connectionState != ConnectionState.done) {
          return const Center(
            child: CircularProgressIndicator(),
          );
        }

        if (snapshot.hasError) {
          return const Center(
            child: Text('Failed to load existing stickers'),
          );
        }

        final images = snapshot.data ?? <StickerImage>[];

        return ScrollConfiguration(
          behavior: const ScrollBehavior().copyWith(
            dragDevices: {
              PointerDeviceKind.touch,
              PointerDeviceKind.mouse,
            },
          ),
          child: Expanded(
            child: GridView.builder(
              shrinkWrap: true,
              gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                crossAxisCount: 5,
                mainAxisSpacing: 10,
                crossAxisSpacing: 10,
              ),
              itemCount: images.length + 1,
              itemBuilder: (BuildContext context, int index) {
                const imageWidth = 250.0;

                if (index == 0) {
                  return Container(
                    padding: const EdgeInsets.all(8),
                    width: imageWidth,
                    child: Column(
                      children: <Widget>[
                        TextField(
                          controller: _imageURLController,
                          decoration: const InputDecoration(
                            labelText: '圖片網址',
                          ),
                        ),
                        Expanded(
                          child: Container(
                            padding: const EdgeInsets.all(16),
                            width: imageWidth,
                            child: ElevatedButton(
                              style: ButtonStyle(
                                shape: WidgetStateProperty.all<
                                    RoundedRectangleBorder>(
                                  RoundedRectangleBorder(
                                    borderRadius: BorderRadius.circular(8.0),
                                  ),
                                ),
                              ),
                              onPressed: () async {
                                await addImage(_imageURLController.text);
                                setState(() {});
                              },
                              child: const Icon(Icons.add),
                            ),
                          ),
                        ),
                      ],
                    ),
                  );
                }

                return Container(
                  padding: const EdgeInsets.all(8),
                  child: DeletableStickerImage(
                    id: images[index - 1].id,
                    url: images[index - 1].url,
                    size: imageWidth,
                    onDeleted: () {
                      setState(() {});
                    },
                  ),
                );
              },
            ),
          ),
        );
      },
    );
  }
}

class StickersView extends StatelessWidget {
  const StickersView({
    super.key,
    required this.stickers,
    required this.toRefresh,
  });

  final List<Sticker> stickers;

  final void Function() toRefresh;

  @override
  Widget build(BuildContext context) {
    return StaggeredGrid.count(
      crossAxisSpacing: 10,
      mainAxisSpacing: 10,
      crossAxisCount: 5,
      children: <Widget>[
        for (final sticker in stickers)
          Center(
            child: StickerOverview(
              stickerID: sticker.id,
              stickerName: sticker.stickerName,
              imageURLs: sticker.images.map((e) => e.url).toList(),
              onDeleted: () {
                toRefresh();
              },
            ),
          ),
      ],
    );
  }
}

class StickerOverview extends StatelessWidget {
  const StickerOverview({
    super.key,
    required this.stickerID,
    required this.stickerName,
    required this.imageURLs,
    required this.onDeleted,
  });

  final int stickerID;

  final String stickerName;

  final List<String> imageURLs;

  final void Function() onDeleted;

  @override
  Widget build(BuildContext context) {
    return MouseRegion(
      cursor: SystemMouseCursors.click,
      child: GestureDetector(
        onTap: () {
          showStickerManagaerDialog(context, stickerName,
                  edit: true, stickerID: stickerID)
              .then(
            (bool? deleted) {
              if (deleted == true) {
                onDeleted();
              }
            },
          );
        },
        child: Card(
          child: Padding(
            padding: const EdgeInsets.fromLTRB(24, 8, 24, 16),
            child: Column(
              children: <Widget>[
                Text(
                  stickerName,
                  style: Theme.of(context).textTheme.titleLarge,
                ),
                const SizedBox(height: 8),
                StickerOverviewImages(imageURLs: imageURLs),
              ],
            ),
          ),
        ),
      ),
    );
  }
}

class StickerOverviewImages extends StatelessWidget {
  const StickerOverviewImages({
    super.key,
    required this.imageURLs,
  });

  final List<String> imageURLs;

  @override
  Widget build(BuildContext context) {
    const width = 250.0;

    if (imageURLs.isEmpty) {
      return const SizedBox(
        width: width,
        child: Center(
          child: Text('No images'),
        ),
      );
    }

    if (imageURLs.length == 1) {
      return SquareNetworkImage(url: imageURLs[0], size: width);
    }

    return NetworkImageGrid(imageURLs: imageURLs, axisCount: 2, size: width);
  }
}

class NetworkImageGrid extends StatelessWidget {
  const NetworkImageGrid({
    super.key,
    required this.imageURLs,
    required this.axisCount,
    required this.size,
  });

  final List<String> imageURLs;
  final int axisCount;
  final double size;

  @override
  Widget build(BuildContext context) {
    const double spacing = 5.0;
    final double imageSize = (size - (axisCount - 1) * spacing) / axisCount;

    return Column(
      mainAxisSize: MainAxisSize.min,
      children: <Widget>[
        for (int i = 0; i < axisCount; i++)
          Column(
            mainAxisSize: MainAxisSize.min,
            children: <Widget>[
              Padding(
                padding: EdgeInsets.only(top: i == 0 ? 0 : spacing),
                child: Row(
                  mainAxisSize: MainAxisSize.min,
                  children: <Widget>[
                    for (int j = 0; j < axisCount; j++)
                      Padding(
                        padding: EdgeInsets.only(left: j == 0 ? 0 : spacing),
                        child: i * axisCount + j < imageURLs.length
                            ? SquareNetworkImage(
                                url: imageURLs[i * axisCount + j],
                                size: imageSize,
                              )
                            : SizedBox.square(dimension: imageSize),
                      )
                  ],
                ),
              ),
            ],
          ),
      ],
    );
  }
}

class DeletableStickerImage extends ConsumerWidget {
  const DeletableStickerImage({
    super.key,
    required this.id,
    required this.url,
    required this.size,
    required this.onDeleted,
  });

  final int id;

  final String url;

  final double size;

  final void Function() onDeleted;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Stack(
      alignment: Alignment.center,
      children: <Widget>[
        SquareNetworkImage(url: url, size: size),
        Align(
          alignment: Alignment.topRight,
          child: IconButton(
            icon: const Icon(Icons.delete),
            color: Colors.red,
            onPressed: () {
              showDialog<bool>(
                context: context,
                builder: (BuildContext context) {
                  return AlertDialog(
                    title: const Text('確定要刪除嗎？'),
                    actions: <Widget>[
                      TextButton(
                        onPressed: () {
                          Navigator.of(context).pop(false);
                        },
                        child: const Text('取消'),
                      ),
                      TextButton(
                        onPressed: () async {
                          final api =
                              ref.read<UIAPIHandler>(apiHandlerProvider);

                          await api.call(context, (api) {
                            return api.deleteStickerImage(id);
                          });

                          if (!context.mounted) {
                            return;
                          }

                          onDeleted();
                          Navigator.of(context).pop(true);
                        },
                        child: const Text('刪除'),
                      ),
                    ],
                  );
                },
              );
            },
          ),
        ),
      ],
    );
  }
}

class SquareNetworkImage extends StatelessWidget {
  const SquareNetworkImage({
    super.key,
    required this.url,
    required this.size,
  });

  final String url;
  final double size;

  @override
  Widget build(BuildContext context) {
    return Image.network(
      url,
      width: size,
      height: size,
      fit: BoxFit.cover,
      loadingBuilder: (BuildContext context, Widget child,
          ImageChunkEvent? loadingProgress) {
        if (loadingProgress == null) {
          return child;
        }

        return Center(
          child: CircularProgressIndicator(
            value: loadingProgress.expectedTotalBytes != null
                ? loadingProgress.cumulativeBytesLoaded /
                    loadingProgress.expectedTotalBytes!
                : null,
          ),
        );
      },
    );
  }
}
