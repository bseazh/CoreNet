-- Enable RLS on all tables
ALTER TABLE profiles ENABLE ROW LEVEL SECURITY;
ALTER TABLE folders ENABLE ROW LEVEL SECURITY;
ALTER TABLE files ENABLE ROW LEVEL SECURITY;
ALTER TABLE file_shares ENABLE ROW LEVEL SECURITY;

-- Profiles policies
CREATE POLICY "Users can view their own profile" ON profiles FOR SELECT USING (auth.uid() = id);
CREATE POLICY "Users can update their own profile" ON profiles FOR UPDATE USING (auth.uid() = id);
CREATE POLICY "Users can insert their own profile" ON profiles FOR INSERT WITH CHECK (auth.uid() = id);
CREATE POLICY "Users can delete their own profile" ON profiles FOR DELETE USING (auth.uid() = id);

-- Folders policies
CREATE POLICY "Users can view their own folders" ON folders FOR SELECT USING (auth.uid() = user_id);
CREATE POLICY "Users can create their own folders" ON folders FOR INSERT WITH CHECK (auth.uid() = user_id);
CREATE POLICY "Users can update their own folders" ON folders FOR UPDATE USING (auth.uid() = user_id);
CREATE POLICY "Users can delete their own folders" ON folders FOR DELETE USING (auth.uid() = user_id);

-- Files policies
CREATE POLICY "Users can view their own files" ON files FOR SELECT USING (auth.uid() = user_id);
CREATE POLICY "Users can upload their own files" ON files FOR INSERT WITH CHECK (auth.uid() = user_id);
CREATE POLICY "Users can update their own files" ON files FOR UPDATE USING (auth.uid() = user_id);
CREATE POLICY "Users can delete their own files" ON files FOR DELETE USING (auth.uid() = user_id);

-- Public file access for shared files
CREATE POLICY "Public can view shared files" ON files FOR SELECT USING (
  EXISTS (
    SELECT 1 FROM file_shares 
    WHERE file_shares.file_id = files.id 
    AND file_shares.is_public = true 
    AND (file_shares.expires_at IS NULL OR file_shares.expires_at > NOW())
  )
);

-- File shares policies
CREATE POLICY "Users can view shares of their files" ON file_shares FOR SELECT USING (
  EXISTS (SELECT 1 FROM files WHERE files.id = file_shares.file_id AND files.user_id = auth.uid())
);
CREATE POLICY "Users can create shares for their files" ON file_shares FOR INSERT WITH CHECK (
  EXISTS (SELECT 1 FROM files WHERE files.id = file_shares.file_id AND files.user_id = auth.uid())
);
CREATE POLICY "Users can update shares of their files" ON file_shares FOR UPDATE USING (
  EXISTS (SELECT 1 FROM files WHERE files.id = file_shares.file_id AND files.user_id = auth.uid())
);
CREATE POLICY "Users can delete shares of their files" ON file_shares FOR DELETE USING (
  EXISTS (SELECT 1 FROM files WHERE files.id = file_shares.file_id AND files.user_id = auth.uid())
);

-- Shared users can view files shared with them
CREATE POLICY "Users can view files shared with them" ON files FOR SELECT USING (
  EXISTS (
    SELECT 1 FROM file_shares 
    WHERE file_shares.file_id = files.id 
    AND file_shares.shared_with = auth.uid()
    AND (file_shares.expires_at IS NULL OR file_shares.expires_at > NOW())
  )
);
