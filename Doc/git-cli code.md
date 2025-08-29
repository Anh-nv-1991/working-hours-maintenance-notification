# --- CẤU HÌNH CƠ BẢN ---
git config --global user.name "Tên của bạn"
git config --global user.email "email@example.com"

# --- KHỞI TẠO / CLONE ---
git init                     # khởi tạo repo mới
git clone <url>              # clone repo từ remote

# --- TRẠNG THÁI / LOG ---
git status                   # xem trạng thái file
git log --oneline --graph    # xem lịch sử commit gọn

# --- STAGE + COMMIT ---
git add <file>               # stage 1 file
git add .                    # stage tất cả
git commit -m "step 5 -devices add working hours"     # commit với message

# --- BRANCH ---
git branch                   # liệt kê branch
git branch <ten-branch>      # tạo branch mới
git checkout <ten-branch>    # chuyển branch
git checkout -b <ten-branch> # tạo + chuyển branch
git merge <ten-branch>       # gộp branch vào branch hiện tại
git branch -d <ten-branch>   # xóa branch đã merge

# --- PUSH / PULL ---
git remote -v                # xem remote
git remote add origin <url>  # thêm remote
git push -u origin main      # push branch main lần đầu
git push                     # push lần sau
git pull                     # cập nhật code mới từ remote

# --- STASH ---
git stash                    # tạm cất thay đổi
git stash pop                # lấy thay đổi ra lại

# --- RESET / REVERT ---
git reset --hard <hash>      # quay lại 1 commit cụ thể (mất thay đổi)
git revert <hash>            # tạo commit đảo ngược 1 commit

# --- TAG ---
git tag v1.0                 # tạo tag
git push origin v1.0         # push tag
