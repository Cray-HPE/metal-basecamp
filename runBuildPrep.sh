if command -v yum > /dev/null; then
  yum install -y \
    skopeo
elif command -v zypper > /dev/null; then
  zypper install -y -f -l \
    skopeo 
else
    echo "Unsupported package manager or package manager not found -- installing nothing"
    exit 1
fi
