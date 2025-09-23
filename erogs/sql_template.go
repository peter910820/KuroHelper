package erogs

import "fmt"

func BuildFuzzySearchCreatorSQL(search string) string {
	return fmt.Sprintf(`
SELECT json_agg(row_to_json(c))
FROM (
    SELECT
        cr.id,
        cr.name,
        cr.twitter_username,
        cr.blog,
        cr.pixiv,
        (
            SELECT json_agg(game_data)
            FROM (
                SELECT
                    g.gamename,
                    g.sellday,
                    g.median,
                    g.count_all,
                    (
                        SELECT json_agg(
                            json_build_object(
                                'shubetu', s2.shubetu,
                                'shubetu_detail', s2.shubetu_detail,
                                'shubetu_detail_name', s2.shubetu_detail_name
                            )
                        )
                        FROM shokushu s2
                        WHERE s2.creater = cr.id
                          AND s2.game = g.id
                    ) AS shokushu
                FROM gamelist g
                WHERE EXISTS (
                    SELECT 1 
                    FROM shokushu s3
                    WHERE s3.creater = cr.id
                      AND s3.game = g.id
                )
                GROUP BY g.id, g.gamename, g.sellday, g.median, g.count_all
            ) AS game_data
        ) AS games
    FROM createrlist cr
    WHERE cr.name ILIKE '%s'
    LIMIT 1
) AS c;`, search)
}
