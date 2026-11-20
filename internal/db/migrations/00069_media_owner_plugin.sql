-- +goose Up

-- Which plugin stored a file (Phase 33.2 review, WR-01).
--
-- A plugin's privacy answer names the media ids of the records it erased or
-- expired, and the host deletes those files. The id is the plugin's claim,
-- and before this column the host could only check that the file belonged to
-- the website: a plugin could name the logo, or an image on a live page, and
-- it went. owner_plugin is the plugin a form attachment was stored for, set
-- when the file is written (files.keep) and never afterwards; '' is a file of
-- the operator's own. media.Store.DeleteOwned refuses every file that is not
-- the calling plugin's.
--
-- The files plugins stored before this column are found where the shipped
-- contact form names them, a "medien_id" in its stored JSON on the same
-- website. A file that is used by more than one plugin's records, or that a
-- plugin of another shape stored, keeps '' and is refused and reported at
-- its erasure: a file left alone and named is the safe side.
ALTER TABLE media ADD COLUMN owner_plugin TEXT NOT NULL DEFAULT '';

UPDATE media SET owner_plugin = (
    SELECT MIN(ps.plugin_id)
    FROM (SELECT plugin_id, website_id, value FROM plugin_store WHERE json_valid(value)) ps,
         json_tree(ps.value) jt
    WHERE ps.website_id = media.website_id
      AND jt.key = 'medien_id' AND jt.type = 'integer' AND jt.atom = media.id
    HAVING COUNT(DISTINCT ps.plugin_id) = 1
)
WHERE EXISTS (
    SELECT 1
    FROM (SELECT plugin_id, website_id, value FROM plugin_store WHERE json_valid(value)) ps,
         json_tree(ps.value) jt
    WHERE ps.website_id = media.website_id
      AND jt.key = 'medien_id' AND jt.type = 'integer' AND jt.atom = media.id
    GROUP BY ps.website_id
    HAVING COUNT(DISTINCT ps.plugin_id) = 1
);

-- +goose Down
ALTER TABLE media DROP COLUMN owner_plugin;
