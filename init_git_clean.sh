#!/usr/bin/env bash
set -euo pipefail

REMOTE_URL="${1:-}"
if [[ -z "${REMOTE_URL}" ]]; then
  echo "Usage: $0 <git-remote-url>"
  exit 1
fi

timestamp="$(date +%Y%m%d-%H%M%S)"
backup_branch="backup-${timestamp}"
backup_tag="init-backup-${timestamp}"

# 1) init repo nếu cần
if [[ ! -d .git ]]; then
  git init
fi

# 2) cấu hình remote origin
if git remote get-url origin >/dev/null 2>&1; then
  git remote set-url origin "${REMOTE_URL}"
else
  git remote add origin "${REMOTE_URL}"
fi

# 3) tạo BACKUP branch từ trạng thái hiện tại (giữ mọi thứ)
git checkout -B "${backup_branch}"

# stage mọi thứ
git add -A

# commit backup (nếu không có thay đổi thì vẫn tạo commit rỗng để lưu mốc)
if git diff --cached --quiet 2>/dev/null; then
  git commit --allow-empty -m "Backup before main push (${timestamp})"
else
  git commit -m "Backup before main push (${timestamp})"
fi

# push backup branch
git push -u origin "${backup_branch}"

# 4) gắn TAG backup và push
git tag "${backup_tag}"
git push origin "${backup_tag}"

# 5) tạo MAIN orphan (không history), commit snapshot sạch và force-push
git checkout --orphan main
git add -A
git commit -m "Init clean repo (${timestamp})"
git push -u origin main --force

echo
echo "✅ Done."
echo " - Backup branch : ${backup_branch}"
echo " - Backup tag    : ${backup_tag}"
echo " - Main (clean)  : main (force-pushed)"
echo
echo "Gợi ý: đặt 'main' làm default branch trên GitHub (Settings → Branches),"
echo "và giữ lại ${backup_branch}/${backup_tag} để hồi phục khi cần."
