import 'package:flutter/material.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

class ChatsPage extends StatefulWidget {
  const ChatsPage({super.key, required this.channel, required this.groupID});
  final WebSocketChannel channel;
  final String groupID;

  @override
  State<ChatsPage> createState() => _ChatsPageState();
}

class _ChatsPageState extends State<ChatsPage> {
  final _controller = TextEditingController();
  final _scrollController = ScrollController();

  final _messages = <String>[];

  @override
  void initState() {
    super.initState();
    _jumpToBottom();
  }

  @override
  void dispose() {
    _scrollController.dispose();
    _controller.dispose();
    widget.channel.sink.close();
    super.dispose();
  }

  void _jumpToBottom() {
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (!_scrollController.hasClients) return;

      final position = _scrollController.position;
      _scrollController.jumpTo(position.maxScrollExtent);
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: SizedBox(
          width: 600.0,
          child: Card(
            child: Padding(
              padding: const EdgeInsets.all(16.0),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    'Group id: ${widget.groupID}',
                    style: const TextStyle(
                      color: Colors.green,
                      fontSize: 16.0,
                      fontWeight: FontWeight.w800,
                    ),
                  ),
                  const SizedBox(height: 32.0),
                  Expanded(
                    child: StreamBuilder(
                      stream: widget.channel.stream,
                      builder: (_, snapshot) {
                        if (widget.channel.closeCode != null) {
                          Navigator.pop(context);
                        }

                        if (snapshot.hasData) {
                          _messages.add(snapshot.data.toString());
                          _jumpToBottom();
                        }

                        return ListView.builder(
                          controller: _scrollController,
                          itemCount: _messages.length,
                          itemBuilder: (_, i) {
                            return Padding(
                              padding: const EdgeInsets.all(8.0),
                              child: Text(_messages[i]),
                            );
                          },
                        );
                      },
                    ),
                  ),
                  const SizedBox(height: 24.0),
                  SizedBox(
                    height: 40.0,
                    child: Row(
                      children: [
                        Expanded(
                          child: TextField(
                            controller: _controller,
                            decoration: const InputDecoration(
                              border: OutlineInputBorder(
                                borderSide: BorderSide(color: Colors.black),
                              ),
                            ),
                            onSubmitted: (_) => _send(),
                          ),
                        ),
                        const SizedBox(width: 16.0),
                        ElevatedButton(
                          style: ElevatedButton.styleFrom(
                            foregroundColor: Colors.white,
                            backgroundColor: Colors.purple,
                          ),
                          onPressed: _send,
                          child: const Text('send'),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }

  void _send() {
    if (_controller.text.isEmpty) return;
    widget.channel.sink.add(_controller.text);
  }
}
