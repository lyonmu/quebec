import React, { Component, ErrorInfo, ReactNode } from 'react';
import { AlertCircle, RefreshCw, Home } from 'lucide-react';

interface Props {
  children: ReactNode;
  fallback?: ReactNode;
}

interface State {
  hasError: boolean;
  error: Error | null;
  errorInfo: ErrorInfo | null;
}

class ErrorBoundary extends Component<Props, State> {
  declare public readonly props: Readonly<Props> & Readonly<{ children?: ReactNode }>;
  public state: Readonly<State>;

  constructor(props: Props) {
    super(props);
    this.state = {
      hasError: false,
      error: null,
      errorInfo: null,
    };
  }

  static getDerivedStateFromError(error: Error): State {
    return {
      hasError: true,
      error,
      errorInfo: null,
    };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    // 记录错误到控制台
    console.error('ErrorBoundary caught an error:', error, errorInfo);
    
    super.setState({
      error,
      errorInfo,
    });

    // 这里可以添加错误上报逻辑，例如发送到错误监控服务
    // reportErrorToService(error, errorInfo);
  }

  handleReset = () => {
    super.setState({
      hasError: false,
      error: null,
      errorInfo: null,
    });
  };

  handleReload = () => {
    window.location.reload();
  };

  handleGoHome = () => {
    window.location.hash = '#/';
    this.handleReset();
  };

  render() {
    if (this.state.hasError) {
      // 如果提供了自定义 fallback，使用它
      if (this.props.fallback) {
        return this.props.fallback;
      }

      // 默认错误 UI
      return (
        <div className="min-h-screen flex items-center justify-center bg-slate-50 dark:bg-slate-950 p-4">
          <div className="max-w-2xl w-full bg-white dark:bg-slate-900 rounded-lg shadow-lg border border-slate-200 dark:border-slate-800 p-6">
            <div className="flex items-center gap-3 mb-4">
              <div className="w-12 h-12 rounded-full bg-red-100 dark:bg-red-900/20 flex items-center justify-center">
                <AlertCircle className="w-6 h-6 text-red-600 dark:text-red-400" />
              </div>
              <div>
                <h1 className="text-xl font-semibold text-slate-900 dark:text-white">
                  应用出现错误
                </h1>
                <p className="text-sm text-slate-500 dark:text-slate-400">
                  很抱歉，应用遇到了意外错误
                </p>
              </div>
            </div>

            {this.state.error && (
              <div className="mb-6 p-4 bg-slate-50 dark:bg-slate-800 rounded-lg border border-slate-200 dark:border-slate-700">
                <p className="text-sm font-mono text-red-600 dark:text-red-400 break-all mb-2">
                  {this.state.error.toString()}
                </p>
                {this.state.errorInfo && (
                  <details className="mt-2">
                    <summary className="text-xs text-slate-500 dark:text-slate-400 cursor-pointer hover:text-slate-700 dark:hover:text-slate-300">
                      查看错误堆栈
                    </summary>
                    <pre className="mt-2 text-xs text-slate-600 dark:text-slate-400 overflow-auto max-h-48 p-2 bg-slate-100 dark:bg-slate-900 rounded">
                      {this.state.errorInfo.componentStack}
                    </pre>
                  </details>
                )}
              </div>
            )}

            <div className="flex flex-wrap gap-3">
              <button
                onClick={this.handleReset}
                className="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors text-sm font-medium"
              >
                <RefreshCw size={16} />
                重试
              </button>
              <button
                onClick={this.handleReload}
                className="flex items-center gap-2 px-4 py-2 bg-slate-600 hover:bg-slate-700 text-white rounded-lg transition-colors text-sm font-medium"
              >
                <RefreshCw size={16} />
                刷新页面
              </button>
              <button
                onClick={this.handleGoHome}
                className="flex items-center gap-2 px-4 py-2 bg-slate-200 dark:bg-slate-700 hover:bg-slate-300 dark:hover:bg-slate-600 text-slate-900 dark:text-white rounded-lg transition-colors text-sm font-medium"
              >
                <Home size={16} />
                返回首页
              </button>
            </div>

            {process.env.NODE_ENV === 'development' && this.state.errorInfo && (
              <div className="mt-6 pt-6 border-t border-slate-200 dark:border-slate-800">
                <p className="text-xs text-slate-500 dark:text-slate-400 mb-2">
                  开发模式：详细错误信息
                </p>
                <pre className="text-xs text-slate-600 dark:text-slate-400 overflow-auto max-h-64 p-3 bg-slate-100 dark:bg-slate-900 rounded border border-slate-200 dark:border-slate-700">
                  {this.state.error.stack}
                </pre>
              </div>
            )}
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}

export default ErrorBoundary;

