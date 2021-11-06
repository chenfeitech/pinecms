window.CURRENT_INDICES = [{
    "title": "Installation",
    "link": "installation",
    "nodes": [{
        "h2": "Installation",
        "h3": "",
        "h4": "",
        "name": "Installation-h2",
        "content": "\n{video} Laracasts provides a free, thorough introduction to Laravel for newcomers to the framework. It's a great place to start your journey.\n"
    }, {
        "h2": "Installation",
        "h3": "Server Requirements",
        "h4": "",
        "name": "ServerRequirements-h3",
        "content": "The Laravel framework has a few system requirements. All of these requirements are satisfied by the Laravel Homestead virtual machine, so it's highly recommended that you use Homestead as your local Laravel development environment.  However, if you are not using Homestead, you will need to make sure your server meets the following requirements:  \n\nPHP >= 7.1.3\nBCMath PHP Extension\nCtype PHP Extension\nJSON PHP Extension\nMbstring PHP Extension\nOpenSSL PHP Extension\nPDO PHP Extension\nTokenizer PHP Extension\nXML PHP Extension\n\n"
    }, {
        "h2": "Installation",
        "h3": "Installing Laravel",
        "h4": "",
        "name": "InstallingLaravel-h3",
        "content": "Laravel utilizes Composer to manage its dependencies. So, before using Laravel, make sure you have Composer installed on your machine."
    }, {
        "h2": "Installation",
        "h3": "Installing Laravel",
        "h4": "Via Laravel Installer",
        "name": "ViaLaravelInstaller-h4",
        "content": "First, download the Laravel installer using Composer:  composer global require laravel/installer  Make sure to place composer's system-wide vendor bin directory in your $PATH so the laravel executable can be located by your system. This directory exists in different locations based on your operating system; however, some common locations include:  \n\nmacOS: $HOME/.composer/vendor/bin\nGNU / Linux Distributions: $HOME/.config/composer/vendor/bin\nWindows: %USERPROFILE%\\AppData\\Roaming\\Composer\\vendor\\bin\n\n  Once installed, the laravel new command will create a fresh Laravel installation in the directory you specify. For instance, laravel new blog will create a directory named blog containing a fresh Laravel installation with all of Laravel's dependencies already installed:  laravel new blog"
    }, {
        "h2": "Installation",
        "h3": "Installing Laravel",
        "h4": "Via Composer Create-Project",
        "name": "ViaComposerCreate-Project-h4",
        "content": "Alternatively, you may also install Laravel by issuing the Composer create-project command in your terminal:  composer create-project --prefer-dist laravel/laravel blog"
    }, {
        "h2": "Installation",
        "h3": "Installing Laravel",
        "h4": "Local Development Server",
        "name": "LocalDevelopmentServer-h4",
        "content": "If you have PHP installed locally and you would like to use PHP's built-in development server to serve your application, you may use the serve Artisan command. This command will start a development server at http://localhost:8000:  php artisan serve  More robust local development options are available via Homestead and Valet."
    }, {
        "h2": "Installation",
        "h3": "Configuration",
        "h4": "Public Directory",
        "name": "PublicDirectory-h4",
        "content": "After installing Laravel, you should configure your web server's document / web root to be the public directory. The index.php in this directory serves as the front controller for all HTTP requests entering your application."
    }, {
        "h2": "Installation",
        "h3": "Configuration",
        "h4": "Configuration Files",
        "name": "ConfigurationFiles-h4",
        "content": "All of the configuration files for the Laravel framework are stored in the config directory. Each option is documented, so feel free to look through the files and get familiar with the options available to you."
    }, {
        "h2": "Installation",
        "h3": "Configuration",
        "h4": "Directory Permissions",
        "name": "DirectoryPermissions-h4",
        "content": "After installing Laravel, you may need to configure some permissions. Directories within the storage and the bootstrap/cache directories should be writable by your web server or Laravel will not run. If you are using the Homestead virtual machine, these permissions should already be set."
    }, {
        "h2": "Installation",
        "h3": "Configuration",
        "h4": "Application Key",
        "name": "ApplicationKey-h4",
        "content": "The next thing you should do after installing Laravel is set your application key to a random string. If you installed Laravel via Composer or the Laravel installer, this key has already been set for you by the php artisan key:generate command.  Typically, this string should be 32 characters long. The key can be set in the .env environment file. If you have not renamed the .env.example file to .env, you should do that now. If the application key is not set, your user sessions and other encrypted data will not be secure!"
    }, {
        "h2": "Installation",
        "h3": "Configuration",
        "h4": "Additional Configuration",
        "name": "AdditionalConfiguration-h4",
        "content": "Laravel needs almost no other configuration out of the box. You are free to get started developing! However, you may wish to review the config/app.php file and its documentation. It contains several options such as timezone and locale that you may wish to change according to your application.  You may also want to configure a few additional components of Laravel, such as:  \n\nCache\nDatabase\nSession\n\n"
    }, {
        "h2": "Web Server Configuration",
        "h3": "Pretty URLs",
        "h4": "Apache",
        "name": "Apache-h4",
        "content": "Laravel includes a public/.htaccess file that is used to provide URLs without the index.php front controller in the path. Before serving Laravel with Apache, be sure to enable the mod_rewrite module so the .htaccess file will be honored by the server.  If the .htaccess file that ships with Laravel does not work with your Apache installation, try this alternative:  Options +FollowSymLinks -Indexes\nRewriteEngine On\n\nRewriteCond %{HTTP:Authorization} .\nRewriteRule .* - [E=HTTP_AUTHORIZATION:%{HTTP:Authorization}]\n\nRewriteCond %{REQUEST_FILENAME} !-d\nRewriteCond %{REQUEST_FILENAME} !-f\nRewriteRule ^ index.php [L]"
    }, {
        "h2": "Web Server Configuration",
        "h3": "Pretty URLs",
        "h4": "Nginx",
        "name": "Nginx-h4",
        "content": "If you are using Nginx, the following directive in your site configuration will direct all requests to the index.php front controller:  location / {\n    try_files $uri $uri/ /index.php?$query_string;\n}  When using Homestead or Valet, pretty URLs will be automatically configured."
    }]
}]