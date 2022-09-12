create table prime_ministers_tbl(
    id INT NOT NULL AUTO_INCREMENT,
    pm_first_name VARCHAR(20) NOT NULL,
    pm_last_name VARCHAR(20) NOT NULL,
    from_date DATE,
    to_date DATE,
    terms INT,
    office VARCHAR(20),
    party VARCHAR(20),
    PRIMARY KEY ( pm_id )
);

SHOW COLUMNS FROM pmdb.prime_ministers_tbl

ALTER TABLE prime_ministers_tbl CHANGE pm_id id INT NOT NULL AUTO_INCREMENT;

DELETE FROM prime_ministers_tbl WHERE id=5 ORDER BY id DESC LIMIT 1;

INSERT INTO prime_ministers_tbl(first_name,last_name,from_date,to_date,terms,office,party) VALUES ('Robert','Walpole','1721-04-03','1742-02-11',4,'Kings Lynn','Whig');

alter table prime_ministers_tbl set column office varchar (50);
ALTER TABLE prime_ministers_tbl MODIFY COLUMN office VARCHAR(50);

// albumsByArtist queries for albums that have the specified artist name.
func albumsByArtist(name string) ([]Album, error) {
    // An albums slice to hold data from returned rows.
    var albums []Album

    rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
    if err != nil {
        return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
    }
    defer rows.Close()
    // Loop through rows, using Scan to assign column data to struct fields.
    for rows.Next() {
        var alb Album
        if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
            return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
        }
        albums = append(albums, alb)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
    }
    return albums, nil
}