# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class H2rzoom < Formula
  desc ""
  homepage "https://github.com/alfonsfoubert/h2rzoom"
  version "0.2.3"

  on_macos do
    url "https://github.com/alfonsfoubert/h2rzoom/releases/download/v0.2.3/h2rzoom_0.2.3_darwin_all.tar.gz"
    sha256 "59d318d461b3316b1b3d395cb41e063a7bb2504480b99330333f6977eedbea43"

    def install
      bin.install "h2rzoom"
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/alfonsfoubert/h2rzoom/releases/download/v0.2.3/h2rzoom_0.2.3_linux_arm64.tar.gz"
      sha256 "0f60d525f70280aa08f51a67c808365772796df8af1d59cd5b55b77bb4852f9a"

      def install
        bin.install "h2rzoom"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/alfonsfoubert/h2rzoom/releases/download/v0.2.3/h2rzoom_0.2.3_linux_amd64.tar.gz"
      sha256 "93fd01194568d1e2ac522cb48c529859d67fd6d899e624296fff131b9d59fcb4"

      def install
        bin.install "h2rzoom"
      end
    end
  end
end
