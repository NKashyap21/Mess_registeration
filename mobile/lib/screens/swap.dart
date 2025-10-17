import 'package:flutter/material.dart';
import 'package:flutter_frontend/theme/app_theme.dart';
import 'package:flutter_frontend/widgets/button.dart';
import 'package:flutter_frontend/api/api_service.dart';

class SwapScreen extends StatefulWidget {
  const SwapScreen({super.key});

  @override
  State<SwapScreen> createState() => _SwapScreenState();
}

class _SwapScreenState extends State<SwapScreen> {
  bool isLoading = true;
  List<dynamic> availableSwaps = [];
  String? error;
  Map<String, dynamic>? userInfo;
  Map<String, dynamic>? userSwapRequest;

  final TextEditingController _passwordController = TextEditingController();
  String _requestType = 'public';

  @override
  void initState() {
    super.initState();
    _loadInitialData();
  }

  @override
  void dispose() {
    _passwordController.dispose();
    super.dispose();
  }

  Future<void> _loadInitialData() async {
    setState(() {
      isLoading = true;
      error = null;
    });

    try {
      final userResponse = await ApiService.getUserInfo();
      if (userResponse['error'] != null) {
        setState(() {
          error = userResponse['error'];
          isLoading = false;
        });
        return;
      }

      userInfo = userResponse['data'];

      // Try to get user's existing swap request
      await _loadUserSwapRequest();

      await _loadSwapData();
    } catch (e) {
      setState(() {
        error = 'Failed to load initial data: $e';
        isLoading = false;
      });
    }
  }

  Future<void> _loadUserSwapRequest() async {
    try {
      final response = await ApiService.getUserSwapRequest();
      setState(() {
        if (response['error'] != null) {
          // No existing swap request is okay, user can create one
          userSwapRequest = null;
        } else {
          userSwapRequest = response['data'];
        }
      });
    } catch (e) {
      // Error getting user swap request is not critical
      setState(() {
        userSwapRequest = null;
      });
    }
  }

  Future<void> _loadSwapData() async {
    try {
      final response = await ApiService.getSwaps();
      setState(() {
        if (response['error'] != null) {
          error = response['error'];
        } else {
          availableSwaps = response['data'] ?? [];
          error = null;
        }
        isLoading = false;
      });
    } catch (e) {
      setState(() {
        error = 'Failed to load swap data: $e';
        isLoading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        title: const Text(
          'Mess Swapping',
          style: TextStyle(color: AppColors.white, fontSize: 22),
        ),
        centerTitle: true,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, color: AppColors.white),
          onPressed: () => Navigator.pop(context),
        ),
      ),
      body: Column(
        children: [
          Expanded(
            child: Container(
              margin: const EdgeInsets.all(16.0),
              decoration: BoxDecoration(
                color: AppColors.darkGrey,
                borderRadius: BorderRadius.circular(12.0),
              ),
              padding: const EdgeInsets.all(24.0),
              child: isLoading
                  ? const Center(child: CircularProgressIndicator())
                  : error != null
                  ? _buildErrorState()
                  : SingleChildScrollView(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          userSwapRequest != null
                              ? _buildCurrentSwapRequestSection()
                              : _buildCreateRequestSection(),
                          const SizedBox(height: 40),
                          _buildAvailableSwapsSection(),
                        ],
                      ),
                    ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildErrorState() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Text(
            'Error: $error',
            style: const TextStyle(color: AppColors.orange, fontSize: 16),
            textAlign: TextAlign.center,
          ),
          const SizedBox(height: 16),
          AppButton(
            label: 'Retry',
            onPressed: _loadInitialData,
            backgroundColor: AppColors.orange,
            textColor: AppColors.white,
          ),
        ],
      ),
    );
  }

  Widget _buildCurrentSwapRequestSection() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'Your Current Swap Request',
          style: TextStyle(
            color: AppColors.white,
            fontSize: 20,
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 24),

        Container(
          padding: const EdgeInsets.all(20),
          decoration: BoxDecoration(
            color: AppColors.black,
            borderRadius: BorderRadius.circular(8),
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  _buildInfoField(
                    'Type',
                    userSwapRequest!['type'] ?? 'Unknown',
                  ),
                  _buildInfoField(
                    'Direction',
                    userSwapRequest!['direction'] ?? 'Unknown',
                  ),
                ],
              ),
              const SizedBox(height: 16),

              if (userSwapRequest!['type'] == 'friend') ...[
                _buildInfoField(
                  'Password',
                  userSwapRequest!['password'] ?? 'N/A',
                ),
                const SizedBox(height: 16),
              ],

              _buildInfoField(
                'Status',
                userSwapRequest!['completed'] ? 'Completed' : 'Active',
              ),
              const SizedBox(height: 16),

              _buildInfoField(
                'Created At',
                _formatDate(userSwapRequest!['created_at']),
              ),
              const SizedBox(height: 24),

              Center(
                child: SizedBox(
                  width: 200,
                  child: AppButton(
                    label: 'Delete Request',
                    onPressed: _deleteSwapRequest,
                    backgroundColor: Colors.red,
                    textColor: AppColors.white,
                    borderRadius: BorderRadius.circular(8),
                  ),
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildCreateRequestSection() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'Create Swap Request',
          style: TextStyle(
            color: AppColors.white,
            fontSize: 20,
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 24),

        Container(
          padding: const EdgeInsets.all(20),
          decoration: BoxDecoration(
            color: AppColors.black,
            borderRadius: BorderRadius.circular(8),
          ),
          child: Column(
            children: [
              _buildRequestTypeSelector(),
              const SizedBox(height: 16),

              if (_requestType == 'friend') ...[
                _buildInputField(
                  controller: _passwordController,
                  label: 'Password',
                  hint: 'Enter password for friend swap',
                  obscureText: false,
                ),
                const SizedBox(height: 16),
              ],

              Center(
                child: SizedBox(
                  width: 200,
                  child: AppButton(
                    label: 'Create Request',
                    onPressed: _createSwapRequest,
                    backgroundColor: AppColors.orange,
                    textColor: AppColors.white,
                    borderRadius: BorderRadius.circular(8),
                  ),
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildInputField({
    required TextEditingController controller,
    required String label,
    required String hint,
    bool obscureText = false,
    TextInputType keyboardType = TextInputType.text,
  }) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: const TextStyle(
            color: AppColors.lightGrey,
            fontSize: 14,
            fontWeight: FontWeight.w500,
          ),
        ),
        const SizedBox(height: 8),
        TextField(
          controller: controller,
          obscureText: obscureText,
          keyboardType: keyboardType,
          style: const TextStyle(color: AppColors.white, fontSize: 16),
          decoration: InputDecoration(
            hintText: hint,
            hintStyle: const TextStyle(
              color: AppColors.lightGrey,
              fontSize: 14,
            ),
            enabledBorder: const UnderlineInputBorder(
              borderSide: BorderSide(color: AppColors.lightGrey),
            ),
            focusedBorder: const UnderlineInputBorder(
              borderSide: BorderSide(color: AppColors.white),
            ),
            contentPadding: const EdgeInsets.symmetric(vertical: 8),
          ),
        ),
      ],
    );
  }

  Widget _buildRequestTypeSelector() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'Request Type',
          style: TextStyle(
            color: AppColors.lightGrey,
            fontSize: 14,
            fontWeight: FontWeight.w500,
          ),
        ),
        const SizedBox(height: 12),
        Row(
          children: [
            Expanded(child: _buildRadioOption('public', 'Public')),
            Expanded(child: _buildRadioOption('friend', 'Friend')),
          ],
        ),
      ],
    );
  }

  Widget _buildRadioOption(String value, String label) {
    return GestureDetector(
      onTap: () {
        setState(() {
          _requestType = value;
        });
      },
      child: Row(
        children: [
          Radio<String>(
            value: value,
            groupValue: _requestType,
            onChanged: (value) {
              setState(() {
                _requestType = value!;
              });
            },
            activeColor: AppColors.orange,
          ),
          Text(
            label,
            style: const TextStyle(color: AppColors.white, fontSize: 16),
          ),
        ],
      ),
    );
  }

  Widget _buildAvailableSwapsSection() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'Available Swap Requests',
          style: TextStyle(
            color: AppColors.white,
            fontSize: 20,
            fontWeight: FontWeight.bold,
          ),
        ),
        const SizedBox(height: 24),

        if (availableSwaps.isEmpty)
          Container(
            padding: const EdgeInsets.all(20),
            decoration: BoxDecoration(
              color: AppColors.black,
              borderRadius: BorderRadius.circular(8),
            ),
            child: const Center(
              child: Text(
                'No swap requests available',
                style: TextStyle(color: AppColors.lightGrey, fontSize: 16),
              ),
            ),
          )
        else
          ...availableSwaps.map((swap) => _buildSwapCard(swap)).toList(),
      ],
    );
  }

  Widget _buildSwapCard(Map<String, dynamic> swap) {
    return Container(
      margin: const EdgeInsets.only(bottom: 16),
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: AppColors.black,
        borderRadius: BorderRadius.circular(8),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              _buildInfoField('Name', swap['name'] ?? 'Unknown'),
              _buildInfoField('Type', swap['type'] ?? 'Unknown'),
            ],
          ),
          const SizedBox(height: 16),

          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Expanded(
                child: _buildInfoField('Email', swap['email'] ?? 'Unknown'),
              ),
              const SizedBox(width: 16),
              _buildInfoField('Direction', swap['direction'] ?? 'Unknown'),
            ],
          ),
          const SizedBox(height: 24),

          Center(
            child: SizedBox(
              width: 200,
              child: AppButton(
                label: 'Accept Request',
                onPressed: () => _acceptSwapRequest(swap),
                backgroundColor: AppColors.orange,
                textColor: AppColors.white,
                borderRadius: BorderRadius.circular(8),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildInfoField(String label, String value) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: const TextStyle(
            color: AppColors.lightGrey,
            fontSize: 12,
            fontWeight: FontWeight.w500,
          ),
        ),
        const SizedBox(height: 4),
        Text(
          value,
          style: const TextStyle(
            color: AppColors.white,
            fontSize: 16,
            fontWeight: FontWeight.w400,
          ),
        ),
      ],
    );
  }

  Future<void> _createSwapRequest() async {
    if (userInfo == null) {
      _showMessage('User information not loaded', isError: true);
      return;
    }

    if (_requestType == 'friend' && _passwordController.text.isEmpty) {
      _showMessage('Please enter a password for friend request', isError: true);
      return;
    }

    // Show loading
    setState(() {
      isLoading = true;
    });

    try {
      final response = await ApiService.createSwapRequest(
        type: _requestType,
        password: _requestType == 'friend'
            ? _passwordController.text
            : 'default_password',
      );

      setState(() {
        isLoading = false;
      });

      if (response['error'] != null) {
        _showMessage(response['error'], isError: true);
      } else {
        _showMessage('Swap request created successfully!');
        _clearForm();
        await _loadUserSwapRequest(); // Load the newly created request
        _loadSwapData();
      }
    } catch (e) {
      setState(() {
        isLoading = false;
      });
      _showMessage('Failed to create swap request: $e', isError: true);
    }
  }

  Future<void> _acceptSwapRequest(Map<String, dynamic> swap) async {
    final confirmed = await _showConfirmationDialog(
      'Accept Swap Request',
      'Are you sure you want to accept this swap request from ${swap['name']}?',
      'Accept',
      'Cancel',
    );

    if (!confirmed) return;

    setState(() {
      isLoading = true;
    });

    try {
      final response = await ApiService.acceptSwapRequest(
        type: swap['type'] ?? 'public',
        userId: swap['user_id'] ?? 0,
      );

      setState(() {
        isLoading = false;
      });

      if (response['error'] != null) {
        _showMessage(response['error'], isError: true);
      } else {
        _showMessage('Swap request accepted successfully!');
        _loadSwapData(); // Refresh the data
      }
    } catch (e) {
      setState(() {
        isLoading = false;
      });
      _showMessage('Failed to accept swap request: $e', isError: true);
    }
  }

  Future<bool> _showConfirmationDialog(
    String title,
    String content,
    String confirmText,
    String cancelText,
  ) async {
    final result = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        backgroundColor: AppColors.darkGrey,
        title: Text(title, style: const TextStyle(color: AppColors.white)),
        content: Text(
          content,
          style: const TextStyle(color: AppColors.lightGrey),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: Text(
              cancelText,
              style: const TextStyle(color: AppColors.lightGrey),
            ),
          ),
          TextButton(
            onPressed: () => Navigator.pop(context, true),
            child: Text(
              confirmText,
              style: const TextStyle(color: AppColors.orange),
            ),
          ),
        ],
      ),
    );
    return result ?? false;
  }

  void _clearForm() {
    _passwordController.clear();
    setState(() {
      _requestType = 'public';
    });
  }

  Future<void> _deleteSwapRequest() async {
    final confirmed = await _showConfirmationDialog(
      'Delete Swap Request',
      'Are you sure you want to delete your current swap request?',
      'Delete',
      'Cancel',
    );

    if (!confirmed) return;

    setState(() {
      isLoading = true;
    });

    try {
      final response = await ApiService.deleteSwapRequest();

      setState(() {
        isLoading = false;
      });

      if (response['error'] != null) {
        _showMessage(response['error'], isError: true);
      } else {
        _showMessage('Swap request deleted successfully!');
        setState(() {
          userSwapRequest = null;
        });
        _loadSwapData(); // Refresh the data
      }
    } catch (e) {
      setState(() {
        isLoading = false;
      });
      _showMessage('Failed to delete swap request: $e', isError: true);
    }
  }

  String _formatDate(String? dateString) {
    if (dateString == null) return 'Unknown';

    try {
      final date = DateTime.parse(dateString);
      return '${date.day}/${date.month}/${date.year} ${date.hour}:${date.minute.toString().padLeft(2, '0')}';
    } catch (e) {
      return 'Invalid date';
    }
  }

  void _showMessage(String message, {bool isError = false}) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(message),
        backgroundColor: isError ? AppColors.orange : Colors.green,
        duration: const Duration(seconds: 3),
      ),
    );
  }
}
