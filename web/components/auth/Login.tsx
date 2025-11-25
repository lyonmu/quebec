import React, { useState, useEffect, useRef } from 'react';
import { User, Lock, RefreshCw, Loader2, Eye, EyeOff } from 'lucide-react';
import { useLanguage } from '../../contexts/LanguageContext';
import { loginService } from '../../services/system/loginService';
import { CaptchaData } from '../../types';
import forge from 'node-forge';

interface LoginProps {
  onLogin: () => void;
}

const Login: React.FC<LoginProps> = ({ onLogin }) => {
  const { t } = useLanguage();
  const [isLoading, setIsLoading] = useState(false);
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [captchaCode, setCaptchaCode] = useState('');
  const [captchaData, setCaptchaData] = useState<CaptchaData | null>(null);
  const [captchaLoading, setCaptchaLoading] = useState(false);
  const [error, setError] = useState('');
  
  const initialized = useRef(false);

  const loadCaptcha = async () => {
    setCaptchaLoading(true);
    setError('');
    try {
      const response = await loginService.fetchCaptcha();
      if (response.code === 50000) {
        setCaptchaData(response.data);
      } else {
        setError(response.message || 'Failed to load captcha');
      }
    } catch (err) {
      setError('Failed to load captcha. Please try again.');
      console.error('Captcha error:', err);
    } finally {
      setCaptchaLoading(false);
    }
  };

  useEffect(() => {
    if (!initialized.current) {
      initialized.current = true;
      loadCaptcha();
    }
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (!captchaData) {
      setError('Please load captcha first');
      return;
    }

    if (captchaCode.length !== captchaData.length) {
      setError(`Captcha must be ${captchaData.length} characters`);
      return;
    }

    setIsLoading(true);

    // Hash password with SHA256
    const mdPass = forge.md.sha256.create();
    mdPass.update(password);
    const hashedPassword = mdPass.digest().toHex();

    // Hash username with SHA256
    const mdUser = forge.md.sha256.create();
    mdUser.update(username);
    const hashedUsername = mdUser.digest().toHex();

    console.log('Login attempt:', {
      hashedUsername,
      hashedPassword,
      captchaId: captchaData.id,
      captchaCode,
    });

    try {
      const loginResponse = await loginService.login({
        username: hashedUsername,
        password: hashedPassword,
        captcha: captchaCode,
        captcha_id: captchaData.id
      });

      if (loginResponse.code === 50000) {
        // Store token if provided in response
        if (loginResponse.data.token) {
          localStorage.setItem('x-quebec-token', loginResponse.data.token);
        }
        // Store username
        if (loginResponse.data.username) {
          localStorage.setItem('quebec-username', loginResponse.data.username);
        }
        // Store role name
        if (loginResponse.data.role_name) {
          localStorage.setItem('quebec-role-name', loginResponse.data.role_name);
        }
        onLogin();
      } else {
        setError(loginResponse.message || 'Login failed');
        // Refresh captcha on failure
        loadCaptcha();
        setCaptchaCode('');
      }
    } catch (err: any) {
      console.error('Login error:', err);
      setError(err.message || 'Login failed. Please try again.');
      // Refresh captcha on error
      loadCaptcha();
      setCaptchaCode('');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-white dark:bg-slate-950 flex items-center justify-center p-4 transition-colors duration-300">
      <div className="w-full max-w-md">
        {/* Logo and Title */}
        <div className="text-center mb-8">
          <img 
            src="/quebec.png" 
            alt="Quebec Logo" 
            className="w-40 h-20 mx-auto mb-6 rounded-2xl shadow-md"
          />
          <h1 className="text-2xl font-normal text-slate-900 dark:text-white mb-2">
            {t('login.signIn')}
          </h1>
          <p className="text-sm text-slate-600 dark:text-slate-400">
            {t('login.continueToQuebec')}
          </p>
        </div>

        {/* Login Card */}
        <div className="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-lg p-10 shadow-sm">
          <form onSubmit={handleSubmit} className="space-y-5">
            {/* Error Message */}
            {error && (
              <div className="p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-600 dark:text-red-400 text-sm rounded">
                {error}
              </div>
            )}

            {/* Username */}
            <div className="space-y-2">
              <div className="relative">
                <User className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 dark:text-slate-500" size={20} />
                <input 
                  type="text" 
                  className="w-full bg-transparent border border-slate-300 dark:border-slate-700 text-slate-900 dark:text-white pl-11 pr-4 py-3 rounded focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all placeholder:text-slate-400 dark:placeholder:text-slate-600"
                  placeholder={t('login.username') || 'Username'}
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  required
                />
              </div>
            </div>

            {/* Password */}
            <div className="space-y-2">
              <div className="relative">
                <Lock className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 dark:text-slate-500" size={20} />
                <input 
                  type={showPassword ? "text" : "password"}
                  className="w-full bg-transparent border border-slate-300 dark:border-slate-700 text-slate-900 dark:text-white pl-11 pr-10 py-3 rounded focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all placeholder:text-slate-400 dark:placeholder:text-slate-600"
                  placeholder={t('login.password') || 'Password'}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 dark:hover:text-slate-300 transition-colors"
                >
                  {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
                </button>
              </div>
            </div>

            {/* Captcha */}
            <div className="space-y-2">
              <div className="flex gap-2 items-start">
                <div className="flex-1">
                  <input 
                    type="text" 
                    className="w-full bg-transparent border border-slate-300 dark:border-slate-700 text-slate-900 dark:text-white px-4 py-3 rounded focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all placeholder:text-slate-400 dark:placeholder:text-slate-600"
                    placeholder={t('login.enterCaptcha')}
                    value={captchaCode}
                    onChange={(e) => setCaptchaCode(e.target.value.toUpperCase())}
                    maxLength={captchaData?.length || 6}
                    required
                  />
                </div>
                <button
                  type="button"
                  onClick={loadCaptcha}
                  disabled={captchaLoading}
                  className="relative h-[50px] w-32 border border-slate-300 dark:border-slate-700 rounded overflow-hidden bg-slate-50 dark:bg-slate-800 flex items-center justify-center cursor-pointer hover:opacity-80 transition-opacity disabled:cursor-not-allowed"
                  title="Click to refresh captcha"
                >
                  {captchaLoading ? (
                    <Loader2 className="animate-spin text-slate-400" size={20} />
                  ) : captchaData ? (
                    <img 
                      src={captchaData.pictures} 
                      alt="Captcha" 
                      className="h-full w-full object-cover"
                    />
                  ) : (
                    <span className="text-xs text-slate-400">No captcha</span>
                  )}
                </button>
              </div>
            </div>

            {/* Submit Button */}
            <div className="pt-2">
              <button 
                type="submit" 
                disabled={isLoading}
                className="w-full bg-blue-600 hover:bg-blue-700 text-white py-3 rounded font-medium transition-all disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {isLoading ? (
                  <span className="flex items-center justify-center gap-2">
                    <Loader2 size={18} className="animate-spin" />
                    {t('login.signingIn')}
                  </span>
                ) : (
                  t('login.signIn')
                )}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default Login;