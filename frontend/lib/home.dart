import 'package:flutter/material.dart';
import 'package:frontend/chat.dart';
import 'package:http/http.dart' as http;
import 'package:web_socket_channel/web_socket_channel.dart';

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  final _controller = TextEditingController();
  var _canSubmit = false;
  String? _error;

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Form(
        child: Center(
          child: SizedBox(
            width: 500.0,
            child: Card(
              child: Padding(
                padding: const EdgeInsets.all(16.0),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Text(
                      'Join or Create Group',
                      style: Theme.of(context).textTheme.titleLarge,
                    ),
                    const SizedBox(height: 32.0),
                    TextField(
                      controller: _controller,
                      decoration: const InputDecoration(
                        border: OutlineInputBorder(
                          borderSide: BorderSide(color: Colors.black),
                        ),
                      ),
                      onChanged: (val) {
                        setState(() {
                          if (val.isEmpty) {
                            _canSubmit = false;
                          } else {
                            _canSubmit = true;
                          }
                        });
                      },
                    ),
                    const SizedBox(height: 32.0),
                    _Button(
                      text: 'Join',
                      onTap: !_canSubmit ? null : _join,
                    ),
                    const SizedBox(height: 16.0),
                    _Button(
                      text: 'Create',
                      onTap: !_canSubmit ? null : _create,
                    ),
                    if (_error != null) ...[
                      const SizedBox(height: 32.0),
                      Text(
                        _error!,
                        style: const TextStyle(
                          color: Colors.red,
                          fontSize: 13.0,
                        ),
                      ),
                    ],
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }

  static const baseURL = 'http://localhost:8080';
  static const wsBaseURL = 'ws://localhost:8080';

  Future<WebSocketChannel?> _getChannel() async {
    try {
      final channel = WebSocketChannel.connect(
        Uri.parse('$wsBaseURL/groups/${_controller.text}'),
      );
      await channel.ready;
      return channel;
    } catch (e) {
      print(e);
      setState(() {
        _error = 'Could not establish connection';
      });
      return null;
    }
  }

  Future<void> _join() async {
    final channel = await _getChannel();
    if (channel == null) return;

    if (!mounted) return;
    Navigator.of(context).push(
      MaterialPageRoute(builder: (_) {
        return ChatsPage(
          channel: channel,
          groupID: _controller.text,
        );
      }),
    );
  }

  Future<void> _create() async {
    try {
      final url = Uri.parse("$baseURL/create?id=${_controller.text}");
      await http.post(url);
      await _join();
    } catch (e) {
      print(e);
      setState(() {
        _error = 'Error while creating group';
      });
    }
  }
}

class _Button extends StatelessWidget {
  const _Button({required this.text, required this.onTap});
  final String text;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: double.infinity,
      height: 48.0,
      child: ElevatedButton(
        style: ElevatedButton.styleFrom(
          foregroundColor: Colors.white,
          backgroundColor: Colors.purple,
        ),
        onPressed: onTap,
        child: Text(text),
      ),
    );
  }
}
