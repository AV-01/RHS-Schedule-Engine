import json
import re
import os
import csv
from uuid import uuid4

DATA_FILES = {
    "2022-2023": "data/student_data_22-23.json",
    "2023-2024": "data/student_data_23-24.json",
    "2024-2025": "data/student_data_24-25.json",
    "2025-2026": "data/student_data_25-26.json"
}

OUTPUT_DIR = "output"
os.makedirs(OUTPUT_DIR, exist_ok=True)

def parse_name(raw_name):
    """
    Parses 'Last, First M.' or 'Last, First -.' into (last_name, first_name)

    Strips middle initials, we'll deal with those later
    """
    raw_name = raw_name.strip()

    if ',' not in raw_name:
        return raw_name, "", ""
    
    last, rest = raw_name.split(',', 1)
    last = last.strip()
    rest = rest.strip()

    rest = re.sub(r'\s*-\.\s*$', '', rest).strip() # remove dash and dot

    middle_initial = ""
    mi_match = re.search(r'\s+([A-Z])\.\s*$', rest)
    if mi_match:
        middle_initial = mi_match.group(1)
        rest = rest[:mi_match.start()].strip()
    rest = re.sub(r'\.\s*$', '', rest).strip() # remove lone dots

    # # raw_name = re.sub(r'\s*-\+[A-Z].\s*$', '', raw_name).strip()
    # raw_name = re.sub(r'\s*-\.\s*$', '', raw_name).strip() #remove '-.'
    # raw_name = re.sub(r'\s*\.\s*$', '', raw_name).strip() #remove '.'

    # if ',' not in raw_name:
    #     return raw_name, ""
    
    # last, rest = raw_name.split(',', 1)
    
    # rest = rest.strip()
    # rest = re.sub(r'\s+[A-Z]\.$', '', rest).strip()

    return last.strip(), rest.strip(), middle_initial.strip()

def load_year(filepath):
    with open(filepath, 'r', encoding='utf-8') as f:
        # print(json.load(f))
        return json.load(f)

def build_lookup(students):
    """
    dict keyed like: (last_name_lower, first_name_lower) -> student dict
    """
    lookup = {}
    for s in students:
        last, first, _ = parse_name(s['name'])
        key = (last.lower(), first.lower())
        # print(key)
        lookup[key] = s
    return lookup

year_order = list(DATA_FILES.keys())

all_years = {}

for year_name, filepath in DATA_FILES.items():
    if not os.path.exists(filepath):
        print(f"{filepath} not found! Please fix.")
        continue
    all_years[year_name] = load_year(filepath)
    print(f"success loaded all {year_name} data")

lookups = {}

for y, s in all_years.items():
    lookups[y] = build_lookup(s)

def find_student_id(last, first, base_grade, base_year_index):
    """
    Search years for a matching student to fill student_id
    """
    for offset, year_name in enumerate(year_order[base_year_index + 1:], start = 1):
        lookup = lookups.get(year_name, {})
        key = (last.lower(), first.lower())
        if key in lookup:
            candidate = lookup[key]
            expected_grade = base_grade + offset
            if str(candidate.get('grade', '')) == str(expected_grade):
                sid = candidate.get('student_id', '').strip()
                if sid:
                    return sid
    return ''

school_years_rows = []
year_name_to_id = {}
for i, year_name in enumerate(year_order, start = 1):
    if year_name in all_years:
        school_years_rows.append({'id': i, 'name': year_name})
        year_name_to_id[year_name] = i

students_rows = []
schedules_rows = []

student_registry = {}
used_usernames = set()

for year_index, year_name in enumerate(year_order):
    if year_name not in all_years:
        continue
    year_id = year_name_to_id[year_name]

    for student in all_years[year_name]:
        last, first, middle_initial = parse_name(student['name'])
        grade_str = student.get('grade', '0')
        try:
            grade = int(grade_str)
        except ValueError:
            grade = 0

        raw_student_id = student.get('student_id', '').strip()
        registry_key = (last.lower(), first.lower())

        if registry_key not in student_registry:
            # print(registry_key)
            if raw_student_id:
                sid = raw_student_id
            else:
                sid = find_student_id(last, first, grade, year_index)
            
            base_username = f"{first.lower().replace(' ', '')}.{last.lower().replace(' ', '')}"
            username = base_username
            counter = 1
            while username in used_usernames:
                username = f"{base_username}{counter}"
                counter += 1
            used_usernames.add(username)

            student_uuid = str(uuid4())
            students_rows.append({
                'id':student_uuid,
                'student_id': sid,
                'first_name': first,
                'last_name': last,
                'middle_initial': middle_initial,
                'username': username,
            })

            student_registry[registry_key] = student_uuid

        else:
            student_uuid = student_registry[registry_key]
            if raw_student_id:
                for row in students_rows:
                    if row['id'] == student_uuid and not row['student_id']:
                        row['student_id'] = raw_student_id
                        break
            if middle_initial:
                for row in students_rows:
                    if row['id'] == student_uuid and not row['middle_initial']:
                        row['middle_initial'] = middle_initial
                        break
        
        for period_str, class_info in student.get('periods', {}).items():
            try:
                period_num = int(period_str)
            except ValueError:
                continue
        
            schedules_rows.append({
                'id': str(uuid4()),
                'student_uuid': student_uuid,
                'school_year_id': year_id,
                'grade': grade,
                'period': period_num,
                'class_name': class_info.get('class_name', '').strip(),
                'teacher_name': class_info.get('teacher_name', '').strip(),
                'room_num': class_info.get('room_num', '').strip()
            })
    print(f"Processed {year_name}")

with open(f"{OUTPUT_DIR}/school_year.csv", 'w', newline='', encoding='utf-8') as f:
    writer = csv.DictWriter(f, fieldnames = ['id', 'name'])
    writer.writeheader()
    writer.writerows(school_years_rows)

with open(f"{OUTPUT_DIR}/students.csv", 'w', newline='', encoding='utf-8') as f:
    writer = csv.DictWriter(f, fieldnames= ['id', 'student_id', 'first_name', 'middle_initial', 'last_name', 'username'])
    writer.writeheader()
    writer.writerows(students_rows)

with open(f"{OUTPUT_DIR}/schedules.csv", 'w', newline='', encoding='utf-8') as f:
    writer = csv.DictWriter(f, fieldnames=['id', 'student_uuid', 'school_year_id', 'grade', 'period', 'class_name', 'teacher_name', 'room_num'])
    writer.writeheader()
    writer.writerows(schedules_rows)

print("done")