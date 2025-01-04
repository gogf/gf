#!/usr/bin/env bash

# Install gci
echo "Installing gci..."
go install github.com/daixiang0/gci@latest

# Check if the GCI is installed successfully
if ! command -v gci &> /dev/null
then
    echo "gci could not be installed. Please check your Go setup."
    exit 1
fi

# Use GCI to format the code
echo "Running gci to format code..."
gci write \
      --custom-order \
      --skip-generated \
      --skip-vendor \
      -s standard \
      -s blank \
      -s default \
      -s dot \
      -s "prefix(github.com/gogf/gf/v2)" \
      -s "prefix(github.com/gogf/gf/cmd)" \
      -s "prefix(github.com/gogf/gf/contrib)" \
      -s "prefix(github.com/gogf/gf/example)" \
      ./

# Check the code for changes
git diff --name-only --exit-code || if [ $? != 0 ]; then echo "Notice: gci check failed, please gci before pr." && exit 1; fi
echo "gci check pass."

# Add the local domain name to `/etc/hosts`
echo "Adding local domain to /etc/hosts..."
sudo echo "127.0.0.1   local" | sudo tee -a /etc/hosts