FROM archlinux:latest AS aur_publisher

# Install necessary packages for the script
RUN pacman-key init
RUN pacman -Syu --noconfirm \
  && pacman -S --noconfirm \
  git \
  wget \
  openssh \
  && rm -rf /var/cache/pacman/pkg/*

# Set git user configuration
RUN git config --global user.name "AUR PublisherBot" \
  && git config --global user.email "benjamin+aur_release@chausse.xyz"

# Set the working directory to /root
WORKDIR /root

# Copy the update_aur.sh script to /usr/local/bin and set it as the entrypoint
COPY update_aur.sh /usr/local/bin/update_aur.sh
RUN chmod +x /usr/local/bin/update_aur.sh

# Set the script as the ENTRYPOINT (or CMD if you prefer)
ENTRYPOINT ["/usr/local/bin/update_aur.sh"]
